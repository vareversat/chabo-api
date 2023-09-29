package api

import (
	"fmt"
	"net/http"
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
//	@ExternalDocs.description	All data comes from from the Bordeaux Open Data API
//	@ExternalDocs.url			https://opendata.bordeaux-metropole.fr/explore/dataset/previsions_pont_chaban/information/
//	@License.name				MIT
//	@License.url				https://github.com/vareversat/chabo-api/blob/main/LICENSE.md

func GinRouter(mongoClient *mongo.Client) {

	router := gin.Default()
	router.Use(sentrygin.New(sentrygin.Options{}))
	docs.SwaggerInfo.BasePath = "/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	// Set default fallback to the Swagger UI
	router.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/swagger/index.html")
	})
	v1 := router.Group("/v1")
	{
		v1.GET("/healthcheck", controllers.Healthcheck(mongoClient))
		forecasts := v1.Group("/forecasts")
		{
			forecasts.GET("", controllers.GetAllForecats(mongoClient))
			forecasts.GET(":id", controllers.GetForecastByID(mongoClient))
		}
		management := v1.Group("/management")
		{
			management.POST("refresh", controllers.RefreshForcast(mongoClient))
			management.GET("refresh/last", controllers.GetLastRefreshAction(mongoClient))
		}
	}

	if err := router.Run(fmt.Sprintf("%s:%s", os.Getenv("APP_URI"), os.Getenv("APP_PORT"))); err != nil {
		panic(err)
	}

}
