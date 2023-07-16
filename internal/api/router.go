package api

import (
	"fmt"
	"os"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/vareversat/chabo-api/docs"
	"github.com/vareversat/chabo-api/internal/api/controllers"
	"go.mongodb.org/mongo-driver/mongo"
)

//	@Title						Chabo API - The Chaban-Delmas bridge API
//	@Description				You can get every info you need about all the events of the Chaban-Delmas bridge in Bordeaux, France
//	@Version					1.0
//	@BasePath					/v1
//	@Contact.email				dev@vareversat.fr
//	@Produce					json
//	@Scheme						http
//	@ExternalDocs.description	All data are from the Bordeaux Open Data API
//	@ExternalDocs.url			https://opendata.bordeaux-metropole.fr/explore/dataset/previsions_pont_chaban/information/
//	@License.name				MIT
//	@License.url				http://github.com/vareversat/chabo-api/LICENSE.md

func GinRouter(mongoClient *mongo.Client) {

	router := gin.Default()
	router.Use(sentrygin.New(sentrygin.Options{}))
	docs.SwaggerInfo.BasePath = "/v1"
	v1 := router.Group("/v1")
	{
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
		v1.GET("/forecasts", controllers.GetAllForecats(mongoClient))
		v1.GET("/forecasts/:id", controllers.GetForecastByID(mongoClient))
		v1.POST("/manage/refresh", controllers.RefreshForcast(mongoClient))
		v1.GET("/manage/refresh/last", controllers.GetLastRefreshAction(mongoClient))
		v1.GET("/healthcheck", controllers.Healthcheck(mongoClient))

	}

	router.Run(fmt.Sprintf("%s:%s", os.Getenv("APP_URI"), os.Getenv("APP_PORT")))

}
