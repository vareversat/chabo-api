package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
	_ "time/tzdata"

	"github.com/gin-gonic/gin"
	"github.com/vareversat/chabo-api/internal/domains"
)

// Get the int param paramName passed into the request.
// Return an erro if not specified or empty, the value either
func GetIntParams(c *gin.Context, paramName string) (int, error) {

	paramValue, paramValueExists := c.GetQuery(paramName)

	if !paramValueExists {
		return 0, fmt.Errorf("you have to specify the %s in requests params", paramName)
	} else {
		n, _ := strconv.Atoi(paramValue)
		return n, nil
	}

}

// Get the string param paramName passed into the request.
// Return empty string if not specified or empty, the value either
func GetStringParams(c *gin.Context, paramName string) string {

	paramValue, _ := c.GetQuery(paramName)

	return paramValue

}

// Get the timezone passed into the request header.
// Return an error if the header is missing and/or malformated. The value either
func GetTimezoneFromHeader(c *gin.Context) (*time.Location, error) {

	location, err := time.LoadLocation(c.GetHeader("Timezone"))

	if err != nil {
		return nil, errors.New(
			"you don't have specified a valid identifier for the Timezone. Please refer to https://en.wikipedia.org/wiki/List_of_tz_database_time_zones to use a valid one",
		)
	} else {
		return location, nil
	}

}

// Create the metadalinks associated to the response
func ComputeMetadaLinks(itemCount int, limit int, offset int, path string) []interface{} {
	var links []interface{}

	links = append(links, domains.APIResponseSelfLink{Self: domains.APIResponseLink{Link: path}})
	reOffset := regexp.MustCompile(`offset=\d+`)

	if offset+limit < itemCount {
		newOffset := limit + offset
		newLink := reOffset.ReplaceAllString(path, fmt.Sprintf("offset=%d", newOffset))
		links = append(
			links,
			domains.APIResponseNextLink{Self: domains.APIResponseLink{Link: newLink}},
		)
	}

	if offset != 0 && offset < itemCount && offset-limit >= 0 {
		newOffset := offset - limit
		newLink := reOffset.ReplaceAllString(path, fmt.Sprintf("offset=%d", newOffset))
		links = append(
			links,
			domains.APIResponsePreviousLink{Self: domains.APIResponseLink{Link: newLink}},
		)
	}

	return links
}
