package routers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vareversat/chabo-api/internal/api/controllers"
	"github.com/vareversat/chabo-api/internal/domains"
	"github.com/vareversat/chabo-api/internal/repositories"
	"github.com/vareversat/chabo-api/internal/usecases"
	"go.mongodb.org/mongo-driver/mongo"
)

func ForecastsRouter(timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	forecastRepository := repositories.NewForecastRepository(
		db.Collection(domains.ForecastCollection),
	)
	syncRepository := repositories.NewSyncRepository(db.Collection(domains.SyncCollection))
	forecastController := &controllers.ForecastController{
		ForecastUsecase: usecases.NewForecastUsecase(
			forecastRepository,
			syncRepository,
			timeout,
		),
	}

	forecastGroup := group.Group("/forecasts")
	forecastGroup.GET("", forecastController.GetAllForecats())
	forecastGroup.GET(":id", forecastController.GetForecastByID())
	forecastGroup.GET("/today", forecastController.GetTodayForecasts())
	forecastGroup.POST("/sync", forecastController.SyncForecasts())
}
