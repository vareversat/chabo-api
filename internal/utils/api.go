package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/vareversat/chabo-api/internal/domains"
)

var logrus *log.Entry

// InitOpenApi Init the logger
func InitOpenApi(logger *log.Entry) {
	logrus = logger
}

// GetOpenAPIData Get forecasts data from the Open-data API
// Populate the *domains.OpenDataAPIResponse pointer if the data are correct
func GetOpenAPIData(openDataAPIResponse *domains.BordeauxAPIResponse) error {

	data, err := http.Get(os.Getenv("OPENDATA_API_URL"))

	if err != nil {
		logrus.Fatal(err)

		return err
	}

	responseData, err := io.ReadAll(data.Body)
	if err != nil {
		logrus.Fatal(err)

		return err
	}

	err = json.Unmarshal(responseData, openDataAPIResponse)
	if err != nil {
		logrus.Fatal(err)

		return err
	}

	logrus.Info("Open Data fetched with success")

	return nil
}
