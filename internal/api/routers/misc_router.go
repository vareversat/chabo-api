package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/vareversat/chabo-api/internal/api/controllers"
	"go.mongodb.org/mongo-driver/mongo"
)

func MiscRouter(mongoClient *mongo.Client, group *gin.RouterGroup) {
	group.GET("/healthcheck", controllers.Healthcheck(mongoClient))
}
