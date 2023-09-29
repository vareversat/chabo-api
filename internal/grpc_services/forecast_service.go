package grpc_services

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/vareversat/chabo-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	WarningLogger = log.New(os.Stdout, "WARNING: ", log.LUTC|log.Ltime|log.Lshortfile)
	InfoLogger    = log.New(os.Stdout, "INFO: ", log.LUTC|log.Ltime|log.Lshortfile)
	ErrorLogger   = log.New(os.Stdout, "ERROR: ", log.LUTC|log.Ltime|log.Lshortfile)
)

type ForecastService struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewForecastService(collection *mongo.Collection, ctx context.Context) *ForecastService {
	return &ForecastService{collection, ctx}
}

// Get all forecasts matching the filters
func (fs *ForecastService) GetAllForecasts(
	mongoResponse *models.MongoResponse,
	limit int,
	offset int,
	filter bson.D,
) (int, error) {

	cursor, err := getCursor(filter, limit, offset, fs.collection)

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

// Get a forcast by its unique ID
func (fs *ForecastService) GetForecastbyID(forecast *models.Forecast, ID string) error {

	opts := options.FindOne()
	filter := bson.D{{Key: "_id", Value: ID}}

	cursor := fs.collection.FindOne(context.TODO(), filter, opts)

	err := cursor.Decode(&forecast)

	if err != nil {
		return fmt.Errorf("not found")
	}

	return nil
}

// Insert all forecast to refresh the data
// Return an error and wither or not it failed under cooldown (too many request)
func (fs *ForecastService) InsertAllForecasts(forecasts []models.Forecast) (error, bool) {
	interfaceRecords := make([]interface{}, len(forecasts))

	if true {
		InfoLogger.Printf("Refresh is needed !")
		// Transform to generic interface to be usable by ´coll.InsertMany´
		for i := range forecasts {
			interfaceRecords[i] = forecasts[i]
		}

		filter := bson.D{}

		deleteResult, err := fs.collection.DeleteMany(context.TODO(), filter)
		if err != nil {
			ErrorLogger.Printf(err.Error())

			return err, false
		}
		WarningLogger.Printf(
			"(Delete) %d records deleted in %s !\n",
			deleteResult.DeletedCount,
			"ForecastsCollectionName",
		)

		insertResult, err := fs.collection.InsertMany(context.TODO(), interfaceRecords)
		if err != nil {
			ErrorLogger.Printf(err.Error())

			return err, false
		}
		WarningLogger.Printf(
			"(Insert) %d records inserted in %s !\n",
			len(insertResult.InsertedIDs),
			"ForecastsCollectionName",
		)

		return nil, false
	} else {
		WarningLogger.Printf("Refresh is NOT needed !")
		return fmt.Errorf("you cannot refresh too many time"), true
	}

}

func getCursor(filters bson.D, limit int, offset int, coll *mongo.Collection) (*mongo.Cursor, error) {
	var cursor *mongo.Cursor
	var err error

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
