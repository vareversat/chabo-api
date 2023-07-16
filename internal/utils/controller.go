package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vareversat/chabo-api/internal/models"
)

func GetIntParams(c *gin.Context, paramName string) (int, error) {

	paramValue, paramValueExists := c.GetQuery(paramName)

	if !paramValueExists {
		return 0, fmt.Errorf("you have to specify the %s in requests params", paramName)
	} else {
		n, _ := strconv.Atoi(paramValue)
		return n, nil
	}

}

func GetStringParams(c *gin.Context, paramName string) string {

	paramValue, _ := c.GetQuery(paramName)

	return paramValue

}

func GetTimezone(c *gin.Context) (*time.Location, error) {

	location, err := time.LoadLocation(c.GetHeader("Timezone"))

	if err != nil {
		return nil, errors.New("you don't have specified a valid identifier for the Timezone. Please refer to https://en.wikipedia.org/wiki/List_of_tz_database_time_zones to use a valid one")
	} else {
		return location, nil
	}

}

func GetMetadaLinks(itemCount int, limit int, offset int, path string) []interface{} {
	var links []interface{}

	links = append(links, models.OpenAPISelfLink{Self: models.OpenAPILink{Link: path}})
	reOffset := regexp.MustCompile(`offset=\d+`)

	if offset+limit < itemCount {
		newOffset := limit + offset
		newLink := reOffset.ReplaceAllString(path, fmt.Sprintf("offset=%d", newOffset))
		links = append(links, models.OpenAPINextLink{Self: models.OpenAPILink{Link: newLink}})
	}

	if offset != 0 && offset < itemCount && offset-limit >= 0 {
		newOffset := offset - limit
		newLink := reOffset.ReplaceAllString(path, fmt.Sprintf("offset=%d", newOffset))
		links = append(links, models.OpenAPIPreviousLink{Self: models.OpenAPILink{Link: newLink}})
	}

	return links
}
