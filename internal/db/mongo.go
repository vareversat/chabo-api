package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/vareversat/chabo-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ForecastsCollectionName = os.Getenv("MONGO_FORECASTS_COLLECTION_NAME")
	RefreshesCollectionName = os.Getenv("MONGO_REFRESHES_COLLECTION_NAME")
	DatabaseName            = os.Getenv("MONGO_DATABASE_NAME")
	WarningLogger           = log.New(os.Stdout, "WARNING: ", log.LUTC|log.Ltime|log.Lshortfile)
	InfoLogger              = log.New(os.Stdout, "INFO: ", log.LUTC|log.Ltime|log.Lshortfile)
	ErrorLogger             = log.New(os.Stdout, "ERROR: ", log.LUTC|log.Ltime|log.Lshortfile)
)

func ConnectToMongoInstace() *mongo.Client {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	opts := options.Client().
		ApplyURI(os.Getenv("MONGO_DSN")).
		SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	} else {
		// Send a ping command to test the connection
		if err := client.Ping(context.TODO(), nil); err != nil {
			panic(err)
		}
		InfoLogger.Printf("Connected to %s database !\n", DatabaseName)
	}

	return client
}

func GetAllForecasts(
	client *mongo.Client,
	mongoResponse *models.MongoResponse,
	limit int,
	offset int,
	filter bson.D,
) (int, error) {

	cursor, err := getCursor(filter, client, limit, offset)

	for cursor.Next(context.TODO()) {
		if err := cursor.Decode(&mongoResponse); err != nil {
			ErrorLogger.Printf(err.Error())

			return 0, err
		}
	}

	if err != nil {
		ErrorLogger.Printf(err.Error())

		return 0, err
	}

	InfoLogger.Printf("(GetAllForecasts) %d records retrieved !\n", len(mongoResponse.Results))

	return mongoResponse.Count[0].ItemCount, nil

}

func GetForecastbyID(client *mongo.Client, forecast *models.Forecast, ID string) error {

	coll := client.Database(DatabaseName).Collection(ForecastsCollectionName)
	opts := options.FindOne()
	filter := bson.D{{Key: "_id", Value: ID}}

	cursor := coll.FindOne(context.TODO(), filter, opts)

	err := cursor.Decode(&forecast)

	if err != nil {
		return fmt.Errorf("not found")
	}

	return nil
}

// Insert all forecast to refrehs the data
// Return an error and wither or not it failed under cooldown (too many request)
func InsertAllForecasts(client *mongo.Client, forecasts []models.Forecast) (error, bool) {
	interfaceRecords := make([]interface{}, len(forecasts))

	if refreshIsNeeded(client) {
		InfoLogger.Printf("Refresh is needed !")
		// Transform to generic interface to be usable by ´coll.InsertMany´
		for i := range forecasts {
			interfaceRecords[i] = forecasts[i]
		}

		coll := client.Database(DatabaseName).Collection(ForecastsCollectionName)
		filter := bson.D{}

		deleteResult, err := coll.DeleteMany(context.TODO(), filter)
		if err != nil {
			ErrorLogger.Printf(err.Error())

			return err, false
		}
		WarningLogger.Printf(
			"(Delete) %d records deleted in %s !\n",
			deleteResult.DeletedCount,
			ForecastsCollectionName,
		)

		insertResult, err := coll.InsertMany(context.TODO(), interfaceRecords)
		if err != nil {
			ErrorLogger.Printf(err.Error())

			return err, false
		}
		WarningLogger.Printf(
			"(Insert) %d records inserted in %s !\n",
			len(insertResult.InsertedIDs),
			ForecastsCollectionName,
		)

		return nil, false
	} else {
		WarningLogger.Printf("Refresh is NOT needed !")
		return fmt.Errorf("you cannot refresh too many time"), true
	}

}

func InsertRefresh(client *mongo.Client, refresh models.Refresh) error {

	coll := client.Database(DatabaseName).Collection(RefreshesCollectionName)

	_, err := coll.InsertOne(context.TODO(), refresh)
	if err != nil {
		ErrorLogger.Printf(err.Error())

		return err
	}
	WarningLogger.Printf("(Insert) 1 record inserted in %s !\n", RefreshesCollectionName)
	return nil

}

func GetLastRefresh(client *mongo.Client, refresh *models.Refresh) error {

	coll := client.Database(DatabaseName).Collection(RefreshesCollectionName)

	opts := options.FindOne().SetSort(bson.D{{Key: "timestamp", Value: -1}})
	cursor := coll.FindOne(context.TODO(), bson.D{}, opts)

	err := cursor.Decode(&refresh)

	if err != nil {
		return fmt.Errorf("not found")
	}

	return nil

}

// Check the MongoDB connection
func Ping(client *mongo.Client) error {

	err := client.Ping(context.TODO(), nil)

	if err != nil {
		ErrorLogger.Printf(err.Error())

		return err
	}

	return nil

}

// Check if it's possible to perform a data refresh
func refreshIsNeeded(client *mongo.Client) bool {

	var lastRefresh models.Refresh

	// Get the last refresh to be sure this is not too early
	err := GetLastRefresh(client, &lastRefresh)

	if err != nil {
		// An error here means that the collection is empty
		return true
	} else {
		currentTime := time.Now()
		diff := currentTime.Sub(lastRefresh.Timestamp)

		cooldown, _ := strconv.Atoi(os.Getenv("REFRESH_COOLDOWN_MINS"))

		return int(diff.Minutes()) >= cooldown
	}

}

func getCursor(filters bson.D, client *mongo.Client, limit int, offset int) (*mongo.Cursor, error) {
	var cursor *mongo.Cursor
	var err error

	coll := client.Database(DatabaseName).Collection(ForecastsCollectionName)

	// Sort results by circulation_closing_date
	sortPipeline := bson.D{
		{Key: "$sort", Value: bson.D{{Key: "circulation_closing_date", Value: 1}}},
	}
	// Apply all computed filters with a $match
	filtersPipeline := bson.D{{Key: "$match", Value: filters}}
	// Format the result to get the total match of the filters even if limit and/or offset change
	constraintsPipeline := bson.D{
		{
			Key: "$facet",
			Value: bson.D{
				{
					Key: "results",
					Value: bson.A{
						bson.D{{Key: "$skip", Value: offset}},
						bson.D{{Key: "$limit", Value: limit}},
					},
				},
				{Key: "count", Value: bson.A{bson.D{{Key: "$count", Value: "itemCount"}}}},
			},
		},
	}

	cursor, err = coll.Aggregate(
		context.TODO(),
		mongo.Pipeline{sortPipeline, filtersPipeline, constraintsPipeline},
	)

	return cursor, err
}
