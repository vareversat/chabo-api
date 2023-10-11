package main

import (
	"os"

	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
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

func init() {
	// Init Logrus
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	mongoLogger := log.WithFields(log.Fields{
		"channel": "mongo",
	})
	openApiLogger := log.WithFields(log.Fields{
		"channel": "open_api",
	})
	forecastLogger := log.WithFields(log.Fields{
		"channel": "forecast",
	})
	mongoClient = db.InitMongoClient(mongoLogger)
	utils.InitOpenApi(openApiLogger)
	utils.InitForecast(forecastLogger)
	if err := utils.GetOpenAPIData(&openDataForecasts); err != nil {
		panic(err)
	}
}

func main() {
	appLogger := log.WithFields(log.Fields{
		"channel": "app",
	})

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              SentryDSN,
		TracesSampleRate: 1.0,
		EnableTracing:    true,
		Environment:      Env,
	})
	if err != nil {
		appLogger.Fatalf("sentry.Init: %s", err)
	}

	appLogger.Infof(
		"Welcome to Chabo API ! Starting the project in " + Env + " mode (Gin " + GinMode + ")",
	)
	utils.ComputeForecasts(&forecasts, openDataForecasts)
	if err, _ := db.InsertAllForecasts(mongoClient, forecasts); err != nil {
		appLogger.Warning(err)
	}

	api.GinRouter(mongoClient)
}
