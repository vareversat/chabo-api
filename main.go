package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	ge "github.com/vareversat/chabo-api/internal/generate"
	"github.com/vareversat/chabo-api/internal/grpc_services"
	"github.com/vareversat/chabo-api/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	openDataForecasts models.OpenDataAPIResponse
	forecasts         models.Forecasts
	mongoClient       *mongo.Client
	SentryDSN         = os.Getenv("SENTRY_DSN")
	Env               = os.Getenv("ENV")
	GinMode           = os.Getenv("GIN_MODE")

	forecastCollection *mongo.Collection
	forecastService    grpc_services.ForecastService
	redreshCollection  *mongo.Collection
	refreshService     grpc_services.RefreshService
)

func init() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              SentryDSN,
		TracesSampleRate: 1.0,
		EnableTracing:    true,
		Environment:      Env,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	log.Default().
		Println("[CHABO-API] Welcome to Chabo API ! Starting the project in " + Env + " mode")

	ctx := context.TODO()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	opts := options.Client().
		ApplyURI(os.Getenv("MONGO_DSN")).
		SetServerAPIOptions(serverAPI).SetTimeout(5 * time.Second)
	mongoClient, err = mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	} else {
		// Send a ping command to test the connection
		if err := mongoClient.Ping(ctx, nil); err != nil {
			panic(err)
		}
	}
	fmt.Println("MongoDB successfully connected...")

	// Collections
	forecastCollection = mongoClient.Database(os.Getenv("MONGO_DATABASE_NAME")).Collection(os.Getenv("MONGO_FORECASTS_COLLECTION_NAME"))
	redreshCollection = mongoClient.Database(os.Getenv("MONGO_DATABASE_NAME")).Collection(os.Getenv("MONGO_REFRESHES_COLLECTION_NAME"))

	// Services
	forecastService = *grpc_services.NewForecastService(forecastCollection, ctx)
	refreshService = *grpc_services.NewRefreshService(redreshCollection, ctx)
}

func main() {
	healtcheckServer := grpc_services.NewHealthcheckServer(mongoClient)

	grpcServer := grpc.NewServer()
	listener, _ := net.Listen("tcp", ":8080")
	ge.RegisterHealthCheckServiceServer(grpcServer, healtcheckServer)
	reflection.Register(grpcServer)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
