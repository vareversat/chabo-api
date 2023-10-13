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

func RefreshRouter(timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	refreshRepository := repositories.NewRefreshRepository(db.Collection(domains.RefreshCollection))
	refreshContoller := &controllers.RefreshController{
		RefreshUsecase: usecases.NewRefreshUsecase(
			refreshRepository,
			timeout,
		),
	}

	refreshGroup := group.Group("/refresh")
	refreshGroup.GET("/last", refreshContoller.GetLastRefresh())
}
