package routers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vareversat/chabo-api/internal/api/controllers"
	"github.com/vareversat/chabo-api/internal/domains"
	"github.com/vareversat/chabo-api/internal/repositories/mongodb"
	"github.com/vareversat/chabo-api/internal/repositories/postgresql"
	"github.com/vareversat/chabo-api/internal/usecases"
	"go.mongodb.org/mongo-driver/mongo"
)

func SyncRouter(timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	syncRepository := mongodb.NewSyncRepository(db.Collection(domains.SyncCollection))
	syncController := &controllers.SyncController{
		SyncUseCase: usecases.NewSyncUseCase(
			syncRepository,
			timeout,
		),
	}

	syncGroup := group.Group("/syncs")
	syncGroup.GET("/last", syncController.GetLastSyncAction())
}

func PostgresSyncRouter(timeout time.Duration, pool *pgxpool.Pool, group *gin.RouterGroup) {
	syncRepository := postgresql.NewSyncRepository(pool)
	syncController := &controllers.SyncController{
		SyncUseCase: usecases.NewSyncUseCase(
			syncRepository,
			timeout,
		),
	}

	syncGroup := group.Group("/syncs")
	syncGroup.GET("/last", syncController.GetLastSyncAction())
}
