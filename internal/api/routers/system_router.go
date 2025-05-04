package routers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vareversat/chabo-api/internal/api/controllers"
	"github.com/vareversat/chabo-api/internal/repositories/mongodb"
	"github.com/vareversat/chabo-api/internal/usecases"
	"go.mongodb.org/mongo-driver/mongo"
)

func SystemRouter(timeout time.Duration, mongoClient *mongo.Client, group *gin.RouterGroup) {
	healthcheckRepository := mongodb.NewHealthCheckRepository(mongoClient)
	systemController := &controllers.SystemController{
		HealthCheckUseCase: usecases.NewHealthCheckUseCase(
			healthcheckRepository,
			timeout,
		),
	}

	systemGroup := group.Group("/system")
	systemGroup.GET("/healthcheck", systemController.Healthcheck())
}
