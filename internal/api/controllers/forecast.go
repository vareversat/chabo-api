package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/vareversat/chabo-api/internal/db"
	"github.com/vareversat/chabo-api/internal/models"
	"github.com/vareversat/chabo-api/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetAllForecats godoc
//
//	@Summary		Get all foracasts
//	@Description	Fetch all existing forecasts
//	@Tags			Forecasts
//	@Accept			json
//	@Produce		json
//	@Success		200			{object}	models.ForecastsResponse{}
//	@Failure		400			{object}	models.ErrorResponse{}	"Some params are missing and/or not properly formatted fror the requests"
//	@Failure		500			{object}	models.ErrorResponse{}	"An error occured on the server side"
//	@Param			from		query		string					false	"The date to filter from (RFC3339)"		Format(date-time)
//	@Param			limit		query		int						true	"Set the limit of the queried results"	Format(int)	default(10)
//	@Param			offset		query		int						true	"Set the offset of the queried results"	Format(int)	default(0)
//	@Param			reason		query		string					false	"The closing reason"					Enums(boat, maintenance)
//	@Param			boat		query		string					false	"The boat name of the event"
//	@Param			maneuver	query		string					false	"The boat maneuver of the event"								Enums(leaving_bordeaux, entering_in_bordeaux)
//	@Param			Timezone	header		string					false	"Timezone to format the date related fields (TZ identifier)"	default(UTC)
//	@Router			/forecasts [get]
func GetAllForecats(mongoClient *mongo.Client) gin.HandlerFunc {
	fn := func(c *gin.Context) {

		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			hub.Scope().SetTag("controller", "getAllForecats")
		}

		// Default filter = get everything
		mongoFilter := bson.D{}
		var mongoResponse models.MongoResponse

		location, locationErr := utils.GetTimezoneFromHeader(c)
		limit, limitErr := utils.GetIntParams(c, "limit")
		offset, offsetErr := utils.GetIntParams(c, "offset")
		from := utils.GetStringParams(c, "from")
		reason := utils.GetStringParams(c, "reason")
		boat := utils.GetStringParams(c, "boat")
		maneuver := utils.GetStringParams(c, "maneuver")

		if from != "" {
			time, err := time.Parse(time.RFC3339, from)
			if err != nil {
				c.JSON(
					http.StatusBadRequest,
					models.ErrorResponse{
						Error: "your 'from' param is not in RFC3339 format. See https://datatracker.ietf.org/doc/html/rfc3339#section-5.8",
					},
				)
				return
			}
			mongoFilter = append(
				mongoFilter,
				bson.E{
					Key:   "circulation_reopening_date",
					Value: bson.D{{Key: "$gte", Value: time}},
				},
			)
		}

		if reason != "" {
			mongoFilter = append(mongoFilter, bson.E{Key: "closing_reason", Value: reason})
		}

		if boat != "" {
			mongoFilter = append(mongoFilter, bson.E{Key: "boats.name", Value: boat})
		}

		if maneuver != "" {
			mongoFilter = append(mongoFilter, bson.E{Key: "boats.maneuver", Value: maneuver})
		}

		if limitErr != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: limitErr.Error()})
			sentry.CaptureException(limitErr)
			return
		}

		if offsetErr != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: offsetErr.Error()})
			sentry.CaptureException(offsetErr)
			return
		}

		if locationErr != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: locationErr.Error()})
			sentry.CaptureException(locationErr)
			return
		}

		if limit == 0 {
			errMessage := "the limit param need to be greater or equal to 1"
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: errMessage})
			sentry.CaptureException(fmt.Errorf(errMessage))
			return
		}

		itemCount, err := db.GetAllForecasts(
			mongoClient,
			&mongoResponse,
			limit,
			offset,
			mongoFilter,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
			sentry.CaptureException(locationErr)
			return
		}

		// Setting the requested timezone
		for index, forecast := range mongoResponse.Results {
			mongoResponse.Results[index].CirculationClosingDate = forecast.CirculationClosingDate.In(
				location,
			)
			mongoResponse.Results[index].CirculationReopeningDate = forecast.CirculationReopeningDate.In(
				location,
			)
			for index2, boat := range forecast.Boats {
				mongoResponse.Results[index].Boats[index2].ApproximativeCrossingDate = boat.ApproximativeCrossingDate.In(
					location,
				)
			}
		}
		links := utils.ComputeMetadaLinks(
			itemCount,
			limit,
			offset,
			fmt.Sprintf("%s/%s", c.Request.URL.Path, c.Request.URL.RawQuery),
		)

		response := models.ForecastsResponse{
			Links:     links,
			Hits:      itemCount,
			Forecasts: mongoResponse.Results,
			Limit:     limit,
			Offset:    offset,
			Timezone:  location.String(),
		}

		c.JSON(http.StatusOK, response)
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
//	@Success		200			{object}	models.ForecastResponse{}
//	@Failure		404			{object}	models.ErrorResponse{}	"The ID does not match any forecast"
//	@Failure		400			{object}	models.ErrorResponse{}	"Some params are missing and/or not properly formatted fror the requests"
//	@Failure		500			{object}	models.ErrorResponse{}	"An error occured on the server side"
//	@Param			id			path		string					true	"The forecast ID"
//	@Param			Timezone	header		string					false	"Timezone to format the date related fields (TZ identifier)"	default(UTC)
//	@Router			/forecasts/{id} [get]
func GetForecastByID(mongoClient *mongo.Client) gin.HandlerFunc {
	fn := func(c *gin.Context) {

		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			hub.Scope().SetTag("controller", "getForecastByID")
		}

		var forecast models.Forecast

		err := db.GetForecastbyID(mongoClient, &forecast, c.Param("id"))

		if err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
			return
		}

		location, locationErr := utils.GetTimezoneFromHeader(c)

		if locationErr != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: locationErr.Error()})
			return
		}

		// Setting the requested timezone
		forecast.CirculationClosingDate = forecast.CirculationClosingDate.In(location)
		forecast.CirculationReopeningDate = forecast.CirculationReopeningDate.In(location)
		for index, boat := range forecast.Boats {
			forecast.Boats[index].ApproximativeCrossingDate = boat.ApproximativeCrossingDate.In(
				location,
			)
		}

		c.JSON(
			http.StatusOK,
			models.ForecastResponse{Forecast: forecast, Timezone: location.String()},
		)
	}

	return gin.HandlerFunc(fn)
}
