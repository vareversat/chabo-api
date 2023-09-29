package grpc_services

import (
	ge "github.com/vareversat/chabo-api/internal/generate"
	"go.mongodb.org/mongo-driver/mongo"
)

type HealthcheckServer struct {
	ge.UnimplementedHealthCheckServiceServer

	client *mongo.Client
}

func NewHealthcheckServer(client *mongo.Client) *HealthcheckServer {

	return &HealthcheckServer{client: client}
}
