package routers

import (
	"fmt"
	"net/http"
	"os"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/vareversat/chabo-api/docs"
	"github.com/vareversat/chabo-api/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

//	@Title						Chabo API - The Chaban-Delmas bridge API
//	@Description				You can get every info you need about all the events of the Chaban-Delmas bridge in Bordeaux, France
//	@Contact.email				dev@vareversat.fr
//	@Produce					json
//	@Scheme						http
//	@ExternalDocs.description	All data comes from the Bordeaux Open Data API
//	@ExternalDocs.url			https://opendata.bordeaux-metropole.fr/explore/dataset/previsions_pont_chaban/information/
//	@License.name				MIT
//	@License.url				https://github.com/vareversat/chabo-api/blob/main/LICENSE.md

func MainRouter(mongoClient *mongo.Client) {

	// Configure Gin web server
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(utils.JsonLoggerMiddleware())
	router.Use(sentrygin.New(sentrygin.Options{}))
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Version = os.Getenv("API_VERSION")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Set default fallback to the Swagger UI
	router.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/swagger/index.html")
	})

	// Initialize routers
	rootRouterGroup := router.Group(docs.SwaggerInfo.BasePath)

	ForecastsRouter(mongoClient, rootRouterGroup)
	ManagementRouter(mongoClient, rootRouterGroup)
	MiscRouter(mongoClient, rootRouterGroup)

	if err := router.Run(fmt.Sprintf("%s:%s", os.Getenv("APP_URI"), os.Getenv("APP_PORT"))); err != nil {
		panic(err)
	}

}
