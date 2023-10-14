package main

import (
	"os"

	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
	"github.com/vareversat/chabo-api/internal/api/routers"
	"github.com/vareversat/chabo-api/internal/repositories"
	"github.com/vareversat/chabo-api/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	SentryDSN         = os.Getenv("SENTRY_DSN")
	Env               = os.Getenv("ENV")
	GinMode           = os.Getenv("GIN_MODE")
	mongoDatabaseName = os.Getenv("MONGO_DATABASE_NAME")
	mongoDatabase     mongo.Database
)

func init() {
	// Init Logrus
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	openApiLogger := log.WithFields(log.Fields{
		"channel": "open_api",
	})
	forecastLogger := log.WithFields(log.Fields{
		"channel": "forecast",
	})
	utils.InitOpenApi(openApiLogger)
	utils.InitForecast(forecastLogger)
	// Init Mongo
	mongoDatabase = *repositories.NewMongoClient().Database(mongoDatabaseName)
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

	routers.MainRouter(mongoDatabase)
}
