package userApi

import (
	"app/cache"
	"app/database"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func getStateCount(c echo.Context) error {

	json_map := make(map[string]string)
	err := json.NewDecoder(c.Request().Body).Decode(&json_map)
	if err != nil {
		return err
	}

	fmt.Println("Getting counts now")

	fmt.Println("Getting StateCode: ", json_map)

	stateCode := GetStateCode(json_map["lat"], json_map["long"])

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

func Main() {

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/count", getStateCount)

	port := os.Getenv("PORT")

	fmt.Println("Listening on", port)

	// Start server
	e.Logger.Fatal(e.Start(":" + port))

}
