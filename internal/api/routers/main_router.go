package routers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/vareversat/chabo-api/docs"
	"github.com/vareversat/chabo-api/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

// Add the variable to be rewritten during build stage via -ldflags
var version = "undefined"

//	@Title						Chabo API - The Chaban-Delmas bridge API
//	@Description				You can get every info you need about all the events of the Chaban-Delmas bridge in Bordeaux, France
//	@Contact.email				dev@vareversat.fr
//	@Produce					json
//	@Scheme						http
//	@ExternalDocs.description	All data comes from the Bordeaux Open Data API
//	@ExternalDocs.url			https://opendata.bordeaux-metropole.fr/explore/dataset/previsions_pont_chaban/information/
//	@License.name				MIT
//	@License.url				https://github.com/vareversat/chabo-api/blob/main/LICENSE.md

func MainRouter(mongoDatabase mongo.Database) {

	// Configure Gin web server
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(utils.JsonLoggerMiddleware())
	router.Use(sentrygin.New(sentrygin.Options{}))
	router.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(swaggerfiles.Handler, ginSwagger.DefaultModelsExpandDepth(-1)),
	)

	// Configure the swagger
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Version = version

	// Set default fallback to the Swagger UI
	router.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/swagger/index.html")
	})

	// Initialize routers
	rootRouterGroup := router.Group(docs.SwaggerInfo.BasePath)

	timeout := time.Duration(30) * time.Second
	ForecastsRouter(timeout, mongoDatabase, rootRouterGroup)
	SyncRouter(timeout, mongoDatabase, rootRouterGroup)
	SystemRouter(timeout, mongoDatabase.Client(), rootRouterGroup)

	// Compute the app address
	// $PORT is automatically injected by Heroku when the app is deployed
	// Use $LOCAL_PORT when $PORT is not defined
	var port string
	var ok bool
	if port, ok = os.LookupEnv("PORT"); !ok {
		port = os.Getenv("LOCAL_PORT")
	}
	appAddr := fmt.Sprintf("%s:%s", os.Getenv("APP_URI"), port)

	if err := router.Run(appAddr); err != nil {
		panic(err)
	}

}
