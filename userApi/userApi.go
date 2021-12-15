package userApi

import (
	"app/cache"
	dbConn "app/database"
	fileutil "app/fileUtil"
	"app/models"
	"strings"

	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.mongodb.org/mongo-driver/bson"
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

		countData = getCovidData(stateCode)
		cache.StoreInRedis(countData)

	}

	countDataIndia := getCovidData("IN")

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

	// fmt.Println(jsonMap["data"])

	stateCode := jsonMap["data"].([]interface{})[0].(map[string]interface{})["region_code"].(string)

	// fmt.Println("Getting State Code:")
	// fmt.Println("Hello---", stateCode)

	return stateCode
}

func getCovidData(stateCode string) models.StateObject {

	fmt.Println("Starting to get data now for ", stateCode)

	var userCollection = dbConn.Db().Database(fileutil.AppConfigProperties["database"]).Collection(fileutil.AppConfigProperties["collection"])

	filterCursor, err := userCollection.Find(context.TODO(), bson.M{"statecode": stateCode})

	fmt.Println("Data is here")

	if err != nil {
		log.Fatal(err)
	}

	var stateInfoBlocks []models.StateObject

	fmt.Println("Initializing StateObject list")

	if err = filterCursor.All(context.TODO(), &stateInfoBlocks); err != nil {
		log.Fatal(err)
	}

	if len(stateInfoBlocks) == 0 {
		var stateObj models.StateObject
		return stateObj
	}

	fmt.Println("Now filtering it!")

	var maxLastTime = stateInfoBlocks[0]

	for i := 1; i < len(stateInfoBlocks); i++ {

		item := stateInfoBlocks[i]

		if maxLastTime.LastUpdated.Before(item.LastUpdated) {
			maxLastTime = item
		}

	}

	return maxLastTime
}

func Main() {

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/count", getStateCount)

	// Start server
	e.Logger.Fatal(e.Start(":9090"))

}
