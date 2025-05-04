package routers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vareversat/chabo-api/internal/api/controllers"
	"github.com/vareversat/chabo-api/internal/domains"
	"github.com/vareversat/chabo-api/internal/repositories/mongodb"
	"github.com/vareversat/chabo-api/internal/repositories/postgresql"
	"github.com/vareversat/chabo-api/internal/usecases"
	"go.mongodb.org/mongo-driver/mongo"
)

func ForecastsRouter(timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	forecastRepository := mongodb.NewForecastRepository(
		db.Collection(domains.ForecastCollection),
	)
	syncRepository := mongodb.NewSyncRepository(db.Collection(domains.SyncCollection))
	forecastController := &controllers.ForecastController{
		ForecastUseCase: usecases.NewForecastUseCase(
			forecastRepository,
			syncRepository,
			timeout,
		),
	}

	forecastGroup := group.Group("/forecasts")
	forecastGroup.GET("", forecastController.GetAllForecasts())
	forecastGroup.GET(":id", forecastController.GetForecastByID())
	forecastGroup.GET("/today", forecastController.GetTodayForecasts())
	forecastGroup.POST("/sync", forecastController.SyncForecasts())
	forecastGroup.GET("/next", forecastController.GetNextForecast())
	forecastGroup.GET("/current", forecastController.GetCurrentForecast())
}

func PostgresForecastsRouter(timeout time.Duration, pool *pgxpool.Pool, group *gin.RouterGroup) {
	forecastRepository := postgresql.NewForecastRepository(pool)
	syncRepository := postgresql.NewSyncRepository(pool)
	forecastController := &controllers.ForecastController{
		ForecastUseCase: usecases.NewForecastUseCase(
			forecastRepository,
			syncRepository,
			timeout,
		),
	}

	forecastGroup := group.Group("/forecasts")
	forecastGroup.GET("", forecastController.GetAllForecasts())
	forecastGroup.GET(":id", forecastController.GetForecastByID())
	forecastGroup.GET("/today", forecastController.GetTodayForecasts())
	forecastGroup.POST("/sync", forecastController.SyncForecasts())
	forecastGroup.GET("/next", forecastController.GetNextForecast())
	forecastGroup.GET("/current", forecastController.GetCurrentForecast())
}
