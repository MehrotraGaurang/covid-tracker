package main

import (
	fileutil "app/fileUtil"
	"app/populateDb"
	"app/userApi"
	"fmt"
)

func main() {
	fmt.Println("Starting Both APIs")

	// Setting Properties
	fileutil.ReadPropertiesFile("properties/properties.txt")

	// fmt.Println(fileutil.AppConfigProperties)

	fmt.Println("Starting populate db")
	go populateDb.Main()

	fmt.Println("Starting User API")
	userApi.Main()

}
