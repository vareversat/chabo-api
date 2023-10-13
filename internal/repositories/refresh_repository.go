package repositories

import (
	"context"

	"github.com/vareversat/chabo-api/internal/domains"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type refreshRepository struct {
	collection *mongo.Collection
}

func NewRefreshRepository(collection *mongo.Collection) domains.RefreshRepository {
	return &refreshRepository{
		collection: collection,
	}
}

func (rR *refreshRepository) InsertOne(ctx context.Context, refresh domains.Refresh) error {
	_, err := rR.collection.InsertOne(ctx, refresh)

	return err

}

func (rR *refreshRepository) GetLast(ctx context.Context, refresh *domains.Refresh) error {
	opts := options.FindOne().SetSort(bson.D{{Key: "timestamp", Value: -1}})
	cursor := rR.collection.FindOne(ctx, bson.D{}, opts)

	return cursor.Decode(&refresh)
}
