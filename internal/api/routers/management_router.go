package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/vareversat/chabo-api/internal/api/controllers"
	"go.mongodb.org/mongo-driver/mongo"
)

func ManagementRouter(mongoClient *mongo.Client, group *gin.RouterGroup) {
	managementGroup := group.Group("/management")
	managementGroup.POST("refresh", controllers.RefreshForcast(mongoClient))
	managementGroup.GET("refresh/last", controllers.GetLastRefreshAction(mongoClient))
}
