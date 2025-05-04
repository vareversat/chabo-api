package mongodb

import (
	"context"

	"github.com/vareversat/chabo-api/internal/domains"
	"go.mongodb.org/mongo-driver/mongo"
)

type healthcheckRepository struct {
	client *mongo.Client
}

func NewHealthCheckRepository(client *mongo.Client) domains.HealthCheckRepository {
	return &healthcheckRepository{
		client: client,
	}
}

func (rR *healthcheckRepository) GetDBHealth(ctx context.Context) error {
	return rR.client.Ping(ctx, nil)
}
