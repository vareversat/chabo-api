package mongodb

import (
	"context"

	"github.com/vareversat/chabo-api/internal/domains"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type syncRepository struct {
	collection *mongo.Collection
}

func NewSyncRepository(collection *mongo.Collection) domains.SyncRepository {
	return &syncRepository{
		collection: collection,
	}
}

func (rR *syncRepository) InsertOne(ctx context.Context, sync domains.Sync) error {
	_, err := rR.collection.InsertOne(ctx, sync)

	return err

}

func (rR *syncRepository) GetLast(ctx context.Context, sync *domains.Sync) error {
	opts := options.FindOne().SetSort(bson.D{{Key: "timestamp", Value: -1}})
	cursor := rR.collection.FindOne(ctx, bson.D{}, opts)

	return cursor.Decode(&sync)
}
