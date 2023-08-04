package controllers

import (
	"net/http"
	"time"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/vareversat/chabo-api/internal/db"
	"github.com/vareversat/chabo-api/internal/models"
	"github.com/vareversat/chabo-api/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

// RefreshForcast godoc
//
//	@Summary		Refresh the data with the ones from the OpenData API
//	@Description	Get, format et populate database with the data from the OpenData API
//	@Tags			Manage
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.Refresh{}
//	@Failure		500	{object}	models.ErrorResponse{}	"An error occured on the server side"
//	@Failure		429	{object}	models.ErrorResponse{}	"Too many attempt to refresh"
//	@Router			/manage/refresh [post]
func RefreshForcast(mongoClient *mongo.Client) gin.HandlerFunc {
	fn := func(c *gin.Context) {

		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			hub.Scope().SetTag("controller", "refreshForcasts")
		}

		var openDataForecasts models.OpenDataAPIResponse
		var forecasts []models.Forecast

		start := time.Now()
		errGet := utils.GetOpenAPIData(&openDataForecasts)
		if errGet != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: errGet.Error()})
			return
		}
		utils.ComputeForecasts(&forecasts, openDataForecasts)

		errInsert, underCooldown := db.InsertAllForecasts(mongoClient, forecasts)
		// An error occured
		if errInsert != nil && !underCooldown {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: errInsert.Error()})
			return
		}

		// The refresh is under cooldown
		if errInsert != nil && underCooldown {
			c.JSON(http.StatusTooManyRequests, models.ErrorResponse{Error: errInsert.Error()})
			return
		}
		elapsed := time.Since(start)

		response := models.Refresh{ItemCount: len(forecasts), Duration: elapsed, Timestamp: start}

		errInsertRefreshProof := db.InsertRefresh(mongoClient, response)
		if errInsertRefreshProof != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: errInsertRefreshProof.Error()})
			return
		}

		c.JSON(http.StatusOK, response)
	}

	return gin.HandlerFunc(fn)
}

// GetLastRefreshAction godoc
//
//	@Summary		Get the last refresh action
//	@Description	Get the last trace of refresh action on POST /manage/refresh
//	@Tags			Manage
//	@Produce		json
//	@Success		200	{object}	models.Refresh{}
//	@Failure		404	{object}	models.ErrorResponse{}	"No previous refresh action exists"
//	@Failure		500	{object}	models.ErrorResponse{}	"An error occured on the server side"
//	@Router			/manage/refresh/last [get]
func GetLastRefreshAction(mongoClient *mongo.Client) gin.HandlerFunc {
	fn := func(c *gin.Context) {

		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			hub.Scope().SetTag("controller", "getLastrefreshAction")
		}

		var refresh models.Refresh

		err := db.GetLastRefresh(mongoClient, &refresh)

		if err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, refresh)
	}

	return gin.HandlerFunc(fn)
}
