package controllers

import (
	"net/http"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/vareversat/chabo-api/internal/domains"
)

type SystemController struct {
	HealthCheckUsecase domains.HealthCheckUsecase
}

// Healthcheck godoc
//
//	@Summary	Get the status of the API
//	@Tags		System
//	@Produce	json
//	@Success	200	{object}	domains.SystemHealthNOK{}	"The api is healthy"
//	@Failure	503	{object}	domains.SystemHealthOK{}	"The api is unhealthy"
//	@Router		/system/healthcheck [get]
func (sC *SystemController) Healthcheck() gin.HandlerFunc {
	fn := func(c *gin.Context) {

		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			hub.Scope().SetTag("controller", "healthcheck")
		}

		err := sC.HealthCheckUsecase.GetHealth(c)

		if err != nil {
			c.JSON(http.StatusServiceUnavailable, domains.SystemHealthNOK{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, domains.SystemHealthOK{Message: "system is running properly"})
	}

	return gin.HandlerFunc(fn)
}
