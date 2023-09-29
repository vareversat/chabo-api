package grpc_services

import (
	"context"
	"log"

	ge "github.com/vareversat/chabo-api/internal/generate"
	"go.mongodb.org/mongo-driver/mongo"
)

// Check the MongoDB connection
func ping(client *mongo.Client) error {

	err := client.Ping(context.TODO(), nil)

	if err != nil {
		ErrorLogger.Printf(err.Error())

		return err
	}

	return nil

}

func (s *HealthcheckServer) CheckHealth(ctx context.Context, req *ge.Empty) (*ge.HealthCheckResponse, error) {

	err := ping(s.client)

	if err == nil {
		return &ge.HealthCheckResponse{
			Status: ge.HealthCheckResponse_SERVING,
			Info:   "the GRPC backend is running properly. Mongo server is reachable",
		}, nil
	}

	return &ge.HealthCheckResponse{
		Status: ge.HealthCheckResponse_NOT_SERVING,
		Info:   err.Error(),
	}, nil
}

func (s *HealthcheckServer) WatchHealth(req *ge.Empty, srv ge.HealthCheckService_WatchHealthServer) error {

	err := ping(s.client)

	if err == nil {
		resp := ge.HealthCheckResponse{
			Status: ge.HealthCheckResponse_SERVING,
			Info:   "the GRPC backend is running properly. Mongo server is reachable",
		}
		if err := srv.Send(&resp); err != nil {
			log.Println("error generating response")
			return nil
		}
	}

	resp := ge.HealthCheckResponse{
		Status: ge.HealthCheckResponse_NOT_SERVING,
		Info:   err.Error(),
	}

	if err := srv.Send(&resp); err != nil {
		log.Println("error generating response")
		return nil
	}

	return nil
}
