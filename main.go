package main

import (
	"app/API"
	fileutil "app/fileUtil"
	"fmt"
)

func main() {
	fmt.Println("Starting Both APIs")

	// Setting Properties
	fileutil.ReadPropertiesFile("properties/properties.txt")

	// fmt.Println(fileutil.AppConfigProperties)

	fmt.Println("Starting populate db")
	go API.StartPopulating()

	fmt.Println("Starting User API")
	API.Main()

}
