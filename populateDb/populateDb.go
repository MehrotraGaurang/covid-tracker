package populateDb

import (
	"app/constants"
	dbConn "app/database"
	fileutil "app/fileUtil"
	"app/models"
	"context"
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

func Main() {
	doEvery(24*time.Hour, populateDb)
}

func populateDb() {

	fmt.Println(constants.StateCodes)

	url := "https://data.covid19india.org/v4/min/data.min.json"

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

	for key, value := range jsonMap {

		switch concreteVal := value.(type) {

		case map[string]interface{}:

			var stateObject models.StateObject

			stateObject.StateCode = key

			// fmt.Println(stateObject.StateCode)

			confirmed := concreteVal["total"].(map[string]interface{})

			stateObject.LastUpdated = time.Now()

			if confirmed["confirmed"] != nil {
				stateObject.TotalNo = confirmed["confirmed"].(float64)

			} else {
				stateObject.TotalNo = 0
			}

			stateObjects = append(stateObjects, stateObject)

		default:
			fmt.Println()
		}

	}

	updateDb(stateObjects)

}

func updateDb(items []models.StateObject) {

	var userCollection = dbConn.Db().Database(fileutil.AppConfigProperties["database"]).Collection(fileutil.AppConfigProperties["collection"])

	totalCount := 0.0

	for _, item := range items {

		totalCount += item.TotalNo

		_, err := userCollection.InsertOne(context.TODO(), item)

		if err != nil {
			log.Fatal(err)
		}

	}

	var indiaObject models.StateObject

	indiaObject.TotalNo = totalCount
	indiaObject.StateCode = "IN"
	indiaObject.LastUpdated = time.Now()

	insertResult, err := userCollection.InsertOne(context.TODO(), indiaObject)

	if err != nil {
		log.Fatal((err))
	}

	fmt.Println("Inserted record for India total count: ", indiaObject.TotalNo, ", ", insertResult)

}
