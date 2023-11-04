package controllers

import (
	"net/http"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/vareversat/chabo-api/internal/domains"
)

type SyncController struct {
	SyncUsecase domains.SyncUsecase
}

// GetLastSyncAction godoc
//
//	@Summary		Get the last sync action
//	@Description	Get the last trace of sync action on POST /forecasts/sync
//	@Tags			Syncs
//	@Produce		json
//	@Success		200	{object}	domains.Sync{}
//	@Failure		404	{object}	domains.APIErrorResponse{}	"No previous sync action exists"
//	@Failure		500	{object}	domains.APIErrorResponse{}	"An error occured on the server side"
//	@Router			/syncs/last [get]
func (mC *SyncController) GetLastSync() gin.HandlerFunc {
	fn := func(c *gin.Context) {

		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			hub.Scope().SetTag("controller", "getLastsyncAction")
		}

		var sync domains.Sync

		customError := mC.SyncUsecase.GetLast(c, &sync)

		if customError != nil {
			c.JSON(
				customError.GetStatusCode(),
				domains.APIErrorResponse{Error: customError.GetErrorMessage()},
			)
			return
		}

		c.JSON(http.StatusOK, sync)
	}

	return gin.HandlerFunc(fn)
}
