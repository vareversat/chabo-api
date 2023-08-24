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
	forecasts         models.Forecasts
	mongoClient       *mongo.Client
	SentryDSN         = os.Getenv("SENTRY_DSN")
	Env               = os.Getenv("ENV")
	GinMode           = os.Getenv("GIN_MODE")
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

	log.Default().
		Println("[CHABO-API] Welcome to Chabo API ! Starting the project in " + Env + " mode (Gin " + GinMode + ")")

	mongoClient = db.ConnectToMongoInstace()
	if err := utils.GetOpenAPIData(&openDataForecasts); err != nil {
		panic(err)
	}
	utils.ComputeForecasts(&forecasts, openDataForecasts)
	if err, _ := db.InsertAllForecasts(mongoClient, forecasts); err != nil {
		panic(err)
	}

	api.GinRouter(mongoClient)
}
