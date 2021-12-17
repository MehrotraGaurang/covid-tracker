package userApi

import (
	"app/cache"
	"app/database"
	fileutil "app/fileUtil"
	"os"
	"strings"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"net/http"
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

	stateCode := getStateCode(json_map["lat"], json_map["long"])

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

func getStateCode(lat string, long string) string {

	url := fileutil.AppConfigProperties["reverse_geo"] + lat + "," + long

	fmt.Println("Getting state code for lat long: " + lat + "," + long)
	fmt.Println(url)

	spaceClient := http.Client{
		Timeout: time.Minute * 2, // Timeout after 2 minutes
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "spacecount-tutorial")

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)

	if readErr != nil {
		log.Fatal(readErr)
	}

	var jsonMap map[string]interface{}

	jsonErr := json.Unmarshal(body, &jsonMap)

	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	stateCode := jsonMap["data"].([]interface{})[0].(map[string]interface{})["region_code"].(string)

	return stateCode
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
