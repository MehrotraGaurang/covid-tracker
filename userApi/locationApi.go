package userApi

import (
	fileutil "app/fileUtil"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func GetStateCode(lat string, long string) string {

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
