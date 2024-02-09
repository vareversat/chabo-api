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
	ForecastUseCase domains.ForecastUseCase
}

// GetAllForecasts godoc
//
//	@Summary		Get all forecasts
//	@Description	Fetch all existing forecasts
//	@Tags			Forecasts
//	@Accept			json
//	@Produce		json
//	@Success		200			{object}	domains.ForecastsResponse{}
//	@Failure		400			{object}	domains.APIErrorResponse{}	"Some params are missing and/or not properly formatted from the requests"
//	@Failure		500			{object}	domains.APIErrorResponse{}	"An error occurred on the server side"
//	@Param			from		query		string						false	"The date to filter from (RFC3339)"		Format(date-time)
//	@Param			limit		query		int							true	"Set the limit of the queried results"	Format(int)	default(10)
//	@Param			offset		query		int							true	"Set the offset of the queried results"	Format(int)	default(0)
//	@Param			reason		query		string						false	"The closing reason"					Enums(boat, maintenance, wine_festival_boats, special_event)
//	@Param			boat		query		string						false	"The boat name of the event"
//	@Param			maneuver	query		string						false	"The boat maneuver of the event"								Enums(leaving_bordeaux, entering_in_bordeaux)
//	@Param			Timezone	header		string						false	"Timezone to format the date related fields (TZ identifier)"	default(UTC)
//	@Router			/forecasts [get]
func (fC *ForecastController) GetAllForecasts() gin.HandlerFunc {
	fn := func(c *gin.Context) {

		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			hub.Scope().SetTag("controller", "GetAllForecasts")
		}

		var forecasts domains.Forecasts
		var totalItemCount int

		location, locationErr := utils.GetTimezoneFromHeader(c)
		limit, limitErr := utils.GetIntParams(c, "limit")
		offset, offsetErr := utils.GetIntParams(c, "offset")
		reason := utils.GetStringParams(c, "reason")
		boat := utils.GetStringParams(c, "boat")
		maneuver := utils.GetStringParams(c, "maneuver")
		parsedTime, timeErr := time.Parse(time.RFC3339, utils.GetStringParams(c, "from"))

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

		customError := fC.ForecastUseCase.GetAllFiltered(
			c,
			location,
			offset,
			limit,
			parsedTime,
			reason,
			maneuver,
			boat,
			&forecasts,
			&totalItemCount,
		)

		if customError != nil {
			c.JSON(
				customError.GetStatusCode(),
				domains.APIErrorResponse{Error: customError.GetErrorMessage()},
			)
			sentry.CaptureException(locationErr)
			return
		}

		links := utils.ComputeMetadataLinks(
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

// GetTodayForecasts godoc
//
//	@Summary		Get the closing schedule for today
//	@Description	Get the closing schedule for today
//	@Tags			Forecasts
//	@Accept			json
//	@Produce		json
//	@Success		200			{object}	domains.ForecastsResponse{}
//	@Failure		400			{object}	domains.APIErrorResponse{}	"Some params are missing and/or not properly formatted for the requests"
//	@Failure		500			{object}	domains.APIErrorResponse{}	"An error occurred on the server side"
//	@Param			limit		query		int							true	"Set the limit of the queried results"							Format(int)	default(10)
//	@Param			offset		query		int							true	"Set the offset of the queried results"							Format(int)	default(0)
//	@Param			Timezone	header		string						false	"Timezone to format the date related fields (TZ identifier)"	default(UTC)
//	@Router			/forecasts/today [get]
func (fC *ForecastController) GetTodayForecasts() gin.HandlerFunc {
	fn := func(c *gin.Context) {

		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			hub.Scope().SetTag("controller", "GetCurrentForecasts")
		}

		var forecasts domains.Forecasts
		var totalItemCount int

		location, locationErr := utils.GetTimezoneFromHeader(c)
		limit, limitErr := utils.GetIntParams(c, "limit")
		offset, offsetErr := utils.GetIntParams(c, "offset")

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

		customError := fC.ForecastUseCase.GetTodayForecasts(
			c,
			&forecasts,
			offset,
			limit,
			location,
			&totalItemCount,
		)

		if customError != nil {
			c.JSON(
				customError.GetStatusCode(),
				domains.APIErrorResponse{Error: customError.GetErrorMessage()},
			)
			sentry.CaptureException(locationErr)
			return
		}

		links := utils.ComputeMetadataLinks(
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

// SyncForecasts godoc
//
//	@Summary		Sync the forecasts with the ones from the OpenData API
//	@Description	Get, format et populate database with the data from the OpenData API
//	@Tags			Forecasts
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	domains.Sync{}
//	@Failure		500	{object}	domains.APIErrorResponse{}	"An error occurred on the server side"
//	@Failure		429	{object}	domains.APIErrorResponse{}	"Too many attempt to sync"
//	@Router			/forecasts/sync [post]
func (fC *ForecastController) SyncForecasts() gin.HandlerFunc {
	fn := func(c *gin.Context) {

		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			hub.Scope().SetTag("controller", "SyncForecasts")
		}

		sync, customError := fC.ForecastUseCase.TryToSyncAll(c)

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

// GetForecastByID godoc
//
//	@Summary		Get a forecast
//	@Description	Fetch a forecast by his unique ID
//	@Tags			Forecasts
//	@Accept			json
//	@Produce		json
//	@Success		200			{object}	domains.ForecastResponse{}
//	@Failure		404			{object}	domains.APIErrorResponse{}	"The ID does not match any forecast"
//	@Failure		400			{object}	domains.APIErrorResponse{}	"Some params are missing and/or not properly formatted from the requests"
//	@Failure		500			{object}	domains.APIErrorResponse{}	"An error occurred on the server side"
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

		customError := fC.ForecastUseCase.GetByID(c, id, &forecast, location)

		if customError != nil {
			c.JSON(
				customError.GetStatusCode(),
				domains.APIErrorResponse{Error: customError.GetErrorMessage()},
			)
			return
		}

		c.JSON(
			http.StatusOK,
			domains.ForecastResponse{Forecast: forecast, Timezone: location.String()},
		)
	}

	return gin.HandlerFunc(fn)
}

// GetCurrentForecast godoc
//
//	@Summary		Fetch the current forecast
//	@Description	Get the current forecast (the bridge is currently closed)
//	@Tags			Forecasts
//	@Accept			json
//	@Produce		json
//	@Success		200			{object}	domains.ForecastResponse{}
//	@Failure		404			{object}	domains.APIErrorResponse{}	"The Chaban-Delmas bridge is currently open"
//	@Failure		400			{object}	domains.APIErrorResponse{}	"Some params are missing and/or not properly formatted from the requests"
//	@Failure		500			{object}	domains.APIErrorResponse{}	"An error occurred on the server side"
//	@Param			Timezone	header		string						false	"Timezone to format the date related fields (TZ identifier)"	default(UTC)
//	@Router			/forecasts/current [get]
func (fC *ForecastController) GetCurrentForecast() gin.HandlerFunc {
	fn := func(c *gin.Context) {

		var forecast domains.Forecast

		location, locationErr := utils.GetTimezoneFromHeader(c)

		if locationErr != nil {
			c.JSON(http.StatusBadRequest, domains.APIErrorResponse{Error: locationErr.Error()})
			return
		}

		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			hub.Scope().SetTag("controller", "GetCurrentForecast")
		}

		err := c.ShouldBind(&forecast)
		if err != nil {
			c.JSON(http.StatusBadRequest, domains.APIErrorResponse{Error: err.Error()})
			return
		}

		customError := fC.ForecastUseCase.GetCurrentForecast(c, &forecast, location)

		if customError != nil {
			c.JSON(
				customError.GetStatusCode(),
				domains.APIErrorResponse{Error: customError.GetErrorMessage()},
			)
			return
		}

		c.JSON(
			http.StatusOK,
			domains.ForecastResponse{Forecast: forecast, Timezone: location.String()},
		)
	}

	return gin.HandlerFunc(fn)
}

// GetNextForecast godoc
//
//	@Summary		Fetch the next forecast
//	@Description	Get the next forecast (= current forecast if the bridge is closed)
//	@Tags			Forecasts
//	@Accept			json
//	@Produce		json
//	@Success		200			{object}	domains.ForecastResponse{}
//	@Failure		404			{object}	domains.APIErrorResponse{}	"The city of Bordeaux has not yet posted the closing times online"
//	@Failure		400			{object}	domains.APIErrorResponse{}	"Some params are missing and/or not properly formatted from the requests"
//	@Failure		500			{object}	domains.APIErrorResponse{}	"An error occurred on the server side"
//	@Param			Timezone	header		string						false	"Timezone to format the date related fields (TZ identifier)"	default(UTC)
//	@Router			/forecasts/next [get]
func (fC *ForecastController) GetNextForecast() gin.HandlerFunc {
	fn := func(c *gin.Context) {

		var forecast domains.Forecast

		location, locationErr := utils.GetTimezoneFromHeader(c)

		if locationErr != nil {
			c.JSON(http.StatusBadRequest, domains.APIErrorResponse{Error: locationErr.Error()})
			return
		}

		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			hub.Scope().SetTag("controller", "GetNextForecast")
		}

		err := c.ShouldBind(&forecast)
		if err != nil {
			c.JSON(http.StatusBadRequest, domains.APIErrorResponse{Error: err.Error()})
			return
		}

		customError := fC.ForecastUseCase.GetNextForecast(c, &forecast, location)

		if customError != nil {
			c.JSON(
				customError.GetStatusCode(),
				domains.APIErrorResponse{Error: customError.GetErrorMessage()},
			)
			return
		}

		c.JSON(
			http.StatusOK,
			domains.ForecastResponse{Forecast: forecast, Timezone: location.String()},
		)
	}

	return gin.HandlerFunc(fn)
}
