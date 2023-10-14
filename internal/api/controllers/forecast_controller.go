package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/vareversat/chabo-api/internal/domains"
	"github.com/vareversat/chabo-api/internal/utils"
)

type ForecastController struct {
	ForecastUsecase domains.ForecastUsecase
}

// GetAllForecats godoc
//
//	@Summary		Get all foracasts
//	@Description	Fetch all existing forecasts
//	@Tags			Forecasts
//	@Accept			json
//	@Produce		json
//	@Success		200			{object}	domains.ForecastsResponse{}
//	@Failure		400			{object}	domains.APIErrorResponse{}	"Some params are missing and/or not properly formatted fror the requests"
//	@Failure		500			{object}	domains.APIErrorResponse{}	"An error occured on the server side"
//	@Param			from		query		string						false	"The date to filter from (RFC3339)"		Format(date-time)
//	@Param			limit		query		int							true	"Set the limit of the queried results"	Format(int)	default(10)
//	@Param			offset		query		int							true	"Set the offset of the queried results"	Format(int)	default(0)
//	@Param			reason		query		string						false	"The closing reason"					Enums(boat, maintenance)
//	@Param			boat		query		string						false	"The boat name of the event"
//	@Param			maneuver	query		string						false	"The boat maneuver of the event"								Enums(leaving_bordeaux, entering_in_bordeaux)
//	@Param			Timezone	header		string						false	"Timezone to format the date related fields (TZ identifier)"	default(UTC)
//	@Router			/forecasts [get]
func (fC *ForecastController) GetAllForecats() gin.HandlerFunc {
	fn := func(c *gin.Context) {

		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			hub.Scope().SetTag("controller", "GetAllForecats")
		}

		var forecasts domains.Forecasts
		var totalItemCount int

		location, locationErr := utils.GetTimezoneFromHeader(c)
		limit, limitErr := utils.GetIntParams(c, "limit")
		offset, offsetErr := utils.GetIntParams(c, "offset")
		reason := utils.GetStringParams(c, "reason")
		boat := utils.GetStringParams(c, "boat")
		maneuver := utils.GetStringParams(c, "maneuver")
		time, timeErr := time.Parse(time.RFC3339, utils.GetStringParams(c, "from"))

		if timeErr != nil && utils.GetStringParams(c, "from") != "" {
			c.JSON(
				http.StatusBadRequest,
				domains.APIErrorResponse{
					Error: "your 'from' param is not in RFC3339 format. See https://datatracker.ietf.org/doc/html/rfc3339#section-5.8",
				},
			)
			sentry.CaptureException(timeErr)
			return
		}

		if limitErr != nil {
			c.JSON(http.StatusBadRequest, domains.APIErrorResponse{Error: limitErr.Error()})
			sentry.CaptureException(limitErr)
			return
		}

		if offsetErr != nil {
			c.JSON(http.StatusBadRequest, domains.APIErrorResponse{Error: offsetErr.Error()})
			sentry.CaptureException(offsetErr)
			return
		}

		if locationErr != nil {
			c.JSON(http.StatusBadRequest, domains.APIErrorResponse{Error: locationErr.Error()})
			sentry.CaptureException(locationErr)
			return
		}

		if limit == 0 {
			errMessage := "the limit param need to be greater or equal to 1"
			c.JSON(http.StatusBadRequest, domains.APIErrorResponse{Error: errMessage})
			sentry.CaptureException(fmt.Errorf(errMessage))
			return
		}

		err := fC.ForecastUsecase.GetAllFiltered(
			c,
			location,
			offset,
			limit,
			time,
			reason,
			maneuver,
			boat,
			&forecasts,
			&totalItemCount,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, domains.APIErrorResponse{Error: err.Error()})
			sentry.CaptureException(locationErr)
			return
		}

		links := utils.ComputeMetadaLinks(
			totalItemCount,
			limit,
			offset,
			fmt.Sprintf("%s/%s", c.Request.URL.Path, c.Request.URL.RawQuery),
		)

		response := domains.ForecastsResponse{
			Links:     links,
			Hits:      totalItemCount,
			Forecasts: forecasts,
			Limit:     limit,
			Offset:    offset,
			Timezone:  location.String(),
		}

		c.JSON(http.StatusOK, response)
	}

	return gin.HandlerFunc(fn)
}

// RefreshForecasts godoc
//
//	@Summary		Refresh the forecasts with the ones from the OpenData API
//	@Description	Get, format et populate database with the data from the OpenData API
//	@Tags			Forecasts
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	domains.Refresh{}
//	@Failure		500	{object}	domains.APIErrorResponse{}	"An error occured on the server side"
//	@Failure		429	{object}	domains.APIErrorResponse{}	"Too many attempt to refresh"
//	@Router			/forecasts/refresh [post]
func (fC *ForecastController) RefreshForecasts() gin.HandlerFunc {
	fn := func(c *gin.Context) {

		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			hub.Scope().SetTag("controller", "RefreshForecasts")
		}

		refresh, err := fC.ForecastUsecase.RefreshAll(c)

		if err != nil {
			c.JSON(http.StatusInternalServerError, domains.APIErrorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, refresh)
	}

	return gin.HandlerFunc(fn)
}

// GetForecastByID godoc
//
//	@Summary		Get a foracast
//	@Description	Fetch a forecast by his unique ID
//	@Tags			Forecasts
//	@Accept			json
//	@Produce		json
//	@Success		200			{object}	domains.ForecastResponse{}
//	@Failure		404			{object}	domains.APIErrorResponse{}	"The ID does not match any forecast"
//	@Failure		400			{object}	domains.APIErrorResponse{}	"Some params are missing and/or not properly formatted fror the requests"
//	@Failure		500			{object}	domains.APIErrorResponse{}	"An error occured on the server side"
//	@Param			id			path		string						true	"The forecast ID"
//	@Param			Timezone	header		string						false	"Timezone to format the date related fields (TZ identifier)"	default(UTC)
//	@Router			/forecasts/{id} [get]
func (fC *ForecastController) GetForecastByID() gin.HandlerFunc {
	fn := func(c *gin.Context) {

		var forecast domains.Forecast
		id := c.Param("id")

		location, locationErr := utils.GetTimezoneFromHeader(c)

		if locationErr != nil {
			c.JSON(http.StatusBadRequest, domains.APIErrorResponse{Error: locationErr.Error()})
			return
		}

		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			hub.Scope().SetTag("controller", "GetForecastByID")
		}

		err := c.ShouldBind(&forecast)
		if err != nil {
			c.JSON(http.StatusBadRequest, domains.APIErrorResponse{Error: err.Error()})
			return
		}

		err = fC.ForecastUsecase.GetByID(c, id, &forecast, location)

		if err != nil {
			c.JSON(http.StatusNotFound, domains.APIErrorResponse{Error: err.Error()})
			return
		}

		c.JSON(
			http.StatusOK,
			domains.ForecastResponse{Forecast: forecast, Timezone: location.String()},
		)
	}

	return gin.HandlerFunc(fn)
}
