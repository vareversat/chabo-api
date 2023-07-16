package main

import (
	"log"
	"os"

	"github.com/getsentry/sentry-go"
	"github.com/vareversat/chabo-api/internal/api"
	"github.com/vareversat/chabo-api/internal/db"
	"github.com/vareversat/chabo-api/internal/models"
	"github.com/vareversat/chabo-api/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	openDataForecasts models.OpenDataAPIResponse
	forecasts         []models.Forecast
	mongoClient       *mongo.Client
	SentryDSN         = os.Getenv("SENTRY_DSN")
	Env               = os.Getenv("ENV")
)

func main() {

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              SentryDSN,
		TracesSampleRate: 1.0,
		EnableTracing:    true,
		Environment:      Env,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	log.Default().Println("[CHABO-API] Welcome to Chabo API ! Starting the project in " + Env + " mode")

	mongoClient = db.ConnectToMongoInstace()
	utils.GetOpenAPIData(&openDataForecasts)
	utils.ComputeForecasts(&forecasts, openDataForecasts)
	db.InsertAllForecasts(mongoClient, forecasts)

	api.GinRouter(mongoClient)
}
