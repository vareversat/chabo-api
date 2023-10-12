package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/vareversat/chabo-api/internal/api/controllers"
	"go.mongodb.org/mongo-driver/mongo"
)

func ForecastsRouter(mongoClient *mongo.Client, group *gin.RouterGroup) {
	forecastGroup := group.Group("/forecasts")
	forecastGroup.GET("", controllers.GetAllForecats(mongoClient))
	forecastGroup.GET(":id", controllers.GetForecastByID(mongoClient))
}
