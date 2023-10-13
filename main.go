package main

import (
	"context"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
	"github.com/vareversat/chabo-api/internal/api/routers"
	"github.com/vareversat/chabo-api/internal/db"
	"github.com/vareversat/chabo-api/internal/domains"
	"github.com/vareversat/chabo-api/internal/models"
	"github.com/vareversat/chabo-api/internal/repositories"
	"github.com/vareversat/chabo-api/internal/usecases"
	"github.com/vareversat/chabo-api/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	openDataForecasts models.OpenDataAPIResponse
	forecasts         domains.Forecasts
	mongoClient       *mongo.Client
	mongoDatabase     mongo.Database
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
	mongoDatabase = *mongoClient.Database(os.Getenv("MONGO_DATABASE_NAME"))
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
	timeout := time.Duration(30) * time.Second
	forecastRepository := repositories.NewForecastRepository(
		mongoDatabase.Collection(domains.ForecastCollection),
	)
	refreshRepository := repositories.NewRefreshRepository(
		mongoDatabase.Collection(domains.RefreshCollection),
	)
	forecastUsecase := usecases.NewForecastUsecase(
		forecastRepository,
		refreshRepository,
		timeout,
	)

	go func() {
		//tick, _ := strconv.Atoi(os.Getenv("REFRESH_TICK_SECONDS"))
		tick := 300000
		for range time.Tick(time.Second * time.Duration(tick)) {
			appLogger.WithFields(log.Fields{
				"kind": "job",
			}).Infof(
				"trying to refresh data",
			)
			if _, err := forecastUsecase.RefreshAll(context.TODO()); err != nil {
				appLogger.Warning(err)
			}
		}
	}()

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

	routers.MainRouter(mongoClient, mongoDatabase)
}
