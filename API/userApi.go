package API

import (
	"app/cache"
	"app/database"
	"app/models"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// Get Count of State and India
// @Summary Get Count of State and India
// @Description Get Count of State and India based on location lat and long provided
// @Accept json
// @Produce json
// @Param Lat header string true "Latitude Required"
// @Param Long header string true "Longitude Required"
// @Success 200 {object} models.StateObject
// @Router /count [get]
func GetStateCount(c echo.Context) error {

	var latLong models.LatLong

	for key, value := range c.Request().Header {
		fmt.Println(key)
		if key == "Lat" {
			latLong.Lat = value[0]
		} else if key == "Long" {
			latLong.Long = value[0]
		}
	}

	fmt.Println("Getting counts now")

	fmt.Println("Getting StateCode: ", latLong)

	stateCode := GetStateCode(latLong.Lat, latLong.Long)

	countData := cache.GetFromRedis(stateCode)

	if strings.Compare(countData.StateCode, "Not_Found") == 0 {

		fmt.Println("Statecode not found, caching it: ", stateCode)

		countData = database.GetCovidData(stateCode)
		cache.StoreInRedis(countData, 30*time.Minute)

	}

	countDataIndia := database.GetCovidData("IN")

	retData := map[string]interface{}{stateCode: countData, "IN": countDataIndia}

	return c.JSON(http.StatusOK, retData)
}
