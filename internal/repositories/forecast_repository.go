package repositories

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vareversat/chabo-api/internal/commons"
	"github.com/vareversat/chabo-api/internal/domains"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type forecastRepository struct {
	collection *mongo.Collection
}

func NewForecastRepository(collection *mongo.Collection) domains.ForecastRepository {
	return &forecastRepository{
		collection: collection,
	}
}

func (fR *forecastRepository) GetByID(
	ctx context.Context,
	id string,
	forecast *domains.Forecast,
) error {
	opts := options.FindOne()
	filter := bson.D{{Key: "_id", Value: id}}

	cursor := fR.collection.FindOne(context.TODO(), filter, opts)

	return cursor.Decode(&forecast)

}

func (fR *forecastRepository) GetAllFiltered(
	ctx context.Context,
	location *time.Location,
	offset int,
	limit int,
	from time.Time,
	reason string,
	maneuver string,
	boat string,
	forecasts *domains.Forecasts,
	totalItemCount *int,
) error {

	var mongoResponse domains.ForecastMongoResponse
	mongoFilter := bson.D{}

	if from.Second() != 0 {
		mongoFilter = append(
			mongoFilter,
			bson.E{
				Key:   "circulation_reopening_date",
				Value: bson.D{{Key: "$gte", Value: from}},
			},
		)
	}

	if reason != "" {
		mongoFilter = append(mongoFilter, bson.E{Key: "closing_reason", Value: reason})
	}

	if boat != "" {
		mongoFilter = append(mongoFilter, bson.E{Key: "boats.name", Value: boat})
	}

	if maneuver != "" {
		mongoFilter = append(mongoFilter, bson.E{Key: "boats.maneuver", Value: maneuver})
	}

	cursor, err := commons.ComputeMongoCursor(ctx, mongoFilter, fR.collection, limit, offset)
	if err != nil {
		return err
	}

	for cursor.Next(ctx) {
		if err := cursor.Decode(&mongoResponse); err != nil {
			logrus.Info(err.Error())
			return err
		}
	}

	*forecasts = mongoResponse.Results
	*totalItemCount = mongoResponse.Count[0].ItemCount

	return err
}

func (fR *forecastRepository) DeleteAll(
	ctx context.Context,
) (int64, error) {
	deleteResult, err := fR.collection.DeleteMany(context.TODO(), bson.D{})

	return deleteResult.DeletedCount, err

}

func (fR *forecastRepository) InsertAll(
	ctx context.Context, forecasts domains.Forecasts,
) (int, error) {
	interfaceRecords := make([]interface{}, len(forecasts))

	for i := range forecasts {
		interfaceRecords[i] = forecasts[i]
	}
	insertResult, err := fR.collection.InsertMany(context.TODO(), interfaceRecords)

	return len(insertResult.InsertedIDs), err

}
