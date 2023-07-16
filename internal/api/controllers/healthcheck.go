package controllers

import (
	"net/http"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/vareversat/chabo-api/internal/db"
	"github.com/vareversat/chabo-api/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

// Healthcheck godoc
//
//	@Summary	Get the status of the API
//	@Tags		Misc
//	@Produce	json
//	@Success	200	{object}	models.OKResponse{}		"The api is healthy"
//	@Failure	503	{object}	models.ErrorResponse{}	"The api is unhealthy"
//	@Router		/healthcheck [get]
func Healthcheck(mongoClient *mongo.Client) gin.HandlerFunc {
	fn := func(c *gin.Context) {

		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			hub.Scope().SetTag("controller", "healthcheck")
		}

		err := db.Ping(mongoClient)

		if err != nil {
			c.JSON(http.StatusServiceUnavailable, models.ErrorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, models.OKResponse{Message: "chabo-api is healthy"})
	}

	return gin.HandlerFunc(fn)
}
