package API

import (
	"app/cache"
	"app/constants"
	"app/database"
	"app/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func doEvery(d time.Duration, f func()) {
	for range time.Tick(d) {
		f()
	}
}

func StartPopulating() {
	fmt.Println("Started Populating Db")
	doEvery(6*time.Hour, populateDb)
}

func populateDb() {

	url := "https://api.rootnet.in/covid19-in/stats/latest"

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

	var stateObjects []models.StateObject

	lastUpdate, err := time.Parse(constants.Layout, jsonMap["lastRefreshed"].(string))

	if err != nil {
		fmt.Println(err)
		lastUpdate = time.Now()
	}

	if isDataUpdatedAtApi(lastUpdate) {
		var indiaObject models.StateObject

		indiaObject.LastUpdated = lastUpdate
		indiaObject.StateCode = "IN"
		indiaObject.StateName = "India"
		indiaObject.TotalNo = jsonMap["data"].(map[string]interface{})["summary"].(map[string]interface{})["total"].(float64)

		stateObjects = append(stateObjects, indiaObject)

		stateInfo := jsonMap["data"].(map[string]interface{})["regional"].([]interface{})

		for _, value := range stateInfo {

			var stateObject models.StateObject

			value := value.(map[string]interface{})

			stateName := value["loc"].(string)

			stateObject.LastUpdated = lastUpdate
			stateObject.StateCode = constants.StateCodes[stateName]
			stateObject.StateName = stateName
			stateObject.TotalNo = value["totalConfirmed"].(float64)

			stateObjects = append(stateObjects, stateObject)
		}

		database.UpdateDb(stateObjects)

	} else {
		fmt.Printf("\nNo new TimeStamp to update DB!")
	}

}

func isDataUpdatedAtApi(lastUpdated time.Time) bool {

	oldTs := cache.GetTs()

	if oldTs.IsZero() || oldTs.Before(lastUpdated) {
		cache.StoreTs(lastUpdated)
		return true
	}
	return false
}
