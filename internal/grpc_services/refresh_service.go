package grpc_services

import (
	"context"
	"fmt"

	"github.com/vareversat/chabo-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RefreshService struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewRefreshService(collection *mongo.Collection, ctx context.Context) *RefreshService {
	return &RefreshService{collection, ctx}
}

func (rs *RefreshService) GetLastRefresh(refresh *models.Refresh) error {

	opts := options.FindOne().SetSort(bson.D{{Key: "timestamp", Value: -1}})
	cursor := rs.collection.FindOne(context.TODO(), bson.D{}, opts)

	err := cursor.Decode(&refresh)

	if err != nil {
		return fmt.Errorf("not found")
	}

	return nil

}
