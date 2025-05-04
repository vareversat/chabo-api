package mongodb

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
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

	cursor := fR.collection.FindOne(ctx, filter, opts)

	return cursor.Decode(&forecast)
}

func (fR *forecastRepository) GetNextForecast(
	ctx context.Context,
	forecast *domains.Forecast,
	now time.Time,
) error {
	opts := options.FindOne()
	filter := bson.D{
		{
			Key:   "circulation_closing_date",
			Value: bson.D{{Key: "$gte", Value: now}},
		},
	}

	cursor := fR.collection.FindOne(ctx, filter, opts)

	return cursor.Decode(&forecast)
}

func (fR *forecastRepository) GetCurrentForecast(
	ctx context.Context,
	forecast *domains.Forecast,
) error {
	mongoFilter := bson.D{}

	mongoFilter = append(
		mongoFilter,
		bson.E{
			Key:   "circulation_closing_date",
			Value: bson.D{{Key: "$lte", Value: time.Now()}},
		},
	)

	mongoFilter = append(
		mongoFilter,
		bson.E{
			Key:   "circulation_opening_date",
			Value: bson.D{{Key: "$gte", Value: time.Now()}},
		},
	)

	cursor, err := computeMongoCursor(ctx, mongoFilter, fR.collection, 1, 0)
	if err != nil {
		return err
	}

	return cursor.Decode(&forecast)
}

func (fR *forecastRepository) GetAllBetweenTwoDates(
	ctx context.Context,
	offset int,
	limit int,
	from time.Time,
	to time.Time,
	forecasts *domains.Forecasts,
	totalItemCount *int,
) error {
	var mongoResponse domains.ForecastMongoResponse
	mongoFilter := bson.D{}

	mongoFilter = append(
		mongoFilter,
		bson.E{
			Key:   "circulation_closing_date",
			Value: bson.D{{Key: "$gte", Value: from}},
		},
	)

	mongoFilter = append(
		mongoFilter,
		bson.E{
			Key:   "circulation_closing_date",
			Value: bson.D{{Key: "$lt", Value: to}},
		},
	)

	cursor, err := computeMongoCursor(ctx, mongoFilter, fR.collection, limit, offset)
	if err != nil {
		return err
	}

	for cursor.Next(ctx) {
		if err := cursor.Decode(&mongoResponse); err != nil {
			logrus.Info(err.Error())
			return err
		}
	}

	// Test if the cursor.Next succeeded to populate the pointer
	if len(mongoResponse.Results) == 0 {
		*forecasts = domains.Forecasts{}
		*totalItemCount = 0

		return nil
	}

	*forecasts = mongoResponse.Results
	*totalItemCount = mongoResponse.Count[0].ItemCount

	return err
}

func (fR *forecastRepository) GetAllFiltered(
	ctx context.Context,
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

	cursor, err := computeMongoCursor(ctx, mongoFilter, fR.collection, limit, offset)
	if err != nil {
		return err
	}

	for cursor.Next(ctx) {
		if err := cursor.Decode(&mongoResponse); err != nil {
			logrus.Info(err.Error())
			return err
		}
	}

	// Test if the cursor.Next succeeded to populate the pointer
	if len(mongoResponse.Results) == 0 {
		*forecasts = domains.Forecasts{}
		*totalItemCount = 0

		return nil
	}

	*forecasts = mongoResponse.Results
	*totalItemCount = mongoResponse.Count[0].ItemCount

	return err
}

func (fR *forecastRepository) DeleteAll(
	ctx context.Context,
) (int64, error) {
	deleteResult, err := fR.collection.DeleteMany(ctx, bson.D{})

	if deleteResult != nil {
		return deleteResult.DeletedCount, nil
	}
	return 0, err
}

func (fR *forecastRepository) InsertAll(
	ctx context.Context, forecasts domains.Forecasts,
) (int, error) {
	interfaceRecords := make([]interface{}, len(forecasts))

	for i := range forecasts {
		interfaceRecords[i] = forecasts[i]
	}
	insertResult, err := fR.collection.InsertMany(ctx, interfaceRecords)

	if insertResult != nil {
		return len(insertResult.InsertedIDs), nil
	}
	return 0, err
}

func computeMongoCursor(
	ctx context.Context,
	filters bson.D,
	collection *mongo.Collection,
	limit int,
	offset int,
) (*mongo.Cursor, error) {
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

	cursor, err = collection.Aggregate(
		ctx,
		mongo.Pipeline{sortPipeline, filtersPipeline, constraintsPipeline},
	)

	return cursor, err
}
