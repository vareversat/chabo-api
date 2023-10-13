package controllers

import (
	"net/http"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/vareversat/chabo-api/internal/domains"
	"github.com/vareversat/chabo-api/internal/models"
)

type RefreshController struct {
	RefreshUsecase domains.RefreshUsecase
}

// GetLastRefreshAction godoc
//
//	@Summary		Get the last refresh action
//	@Description	Get the last trace of refresh action on POST /forecasts/refresh
//	@Tags			Refreshes
//	@Produce		json
//	@Success		200	{object}	domains.Refresh{}
//	@Failure		404	{object}	models.ErrorResponse{}	"No previous refresh action exists"
//	@Failure		500	{object}	models.ErrorResponse{}	"An error occured on the server side"
//	@Router			/refresh/last [get]
func (mC *RefreshController) GetLastRefresh() gin.HandlerFunc {
	fn := func(c *gin.Context) {

		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			hub.Scope().SetTag("controller", "getLastrefreshAction")
		}

		var refresh domains.Refresh

		err := mC.RefreshUsecase.GetLast(c, &refresh)

		if err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, refresh)
	}

	return gin.HandlerFunc(fn)
}
