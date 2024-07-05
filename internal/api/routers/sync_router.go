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

func SyncRouter(timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	syncRepository := repositories.NewSyncRepository(db.Collection(domains.SyncCollection))
	syncController := &controllers.SyncController{
		SyncUseCase: usecases.NewSyncUseCase(
			syncRepository,
			timeout,
		),
	}

	syncGroup := group.Group("/syncs")
	syncGroup.GET("/last", syncController.GetLastSyncAction())
}
