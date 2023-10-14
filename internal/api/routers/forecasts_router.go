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
	refreshRepository := repositories.NewRefreshRepository(db.Collection(domains.RefreshCollection))
	forecastController := &controllers.ForecastController{
		ForecastUsecase: usecases.NewForecastUsecase(
			forecastRepository,
			refreshRepository,
			timeout,
		),
	}

	forecastGroup := group.Group("/forecasts")
	forecastGroup.GET("", forecastController.GetAllForecats())
	forecastGroup.GET(":id", forecastController.GetForecastByID())
	forecastGroup.POST("/refresh", forecastController.RefreshForecasts())
}
