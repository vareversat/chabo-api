package utils

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/vareversat/chabo-api/internal/models"
)

// Get forecasts data from the Opendata API
// Populate the *models.OpenDataAPIResponse pointer if the data are correct
func GetOpenAPIData(openDataAPIResponse *models.OpenDataAPIResponse) error {

	data, err := http.Get(os.Getenv("OPENDATA_API_URL"))

	if err != nil {
		log.Fatal(err)

		return err
	}

	responseData, err := io.ReadAll(data.Body)
	if err != nil {
		log.Fatal(err)

		return err
	}

	err = json.Unmarshal(responseData, openDataAPIResponse)
	if err != nil {
		log.Fatal(err)

		return err
	}

	return nil
}
