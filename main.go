package main

import (
	"app/API"
	_ "app/docs"
	fileutil "app/fileUtil"
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Covid Tracker
// @description Track covid numbers for a state and India
// @version 1.0
// @schemes https
// @host dry-wave-91626.herokuapp.com
// @BasePath /
func main() {

	// Setting Properties
	fileutil.ReadPropertiesFile("properties/properties.txt")

	// fmt.Println(fileutil.AppConfigProperties)

	fmt.Println("Starting populate db")
	go API.StartPopulating()

	fmt.Println("Starting User API")
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/count", API.GetStateCount)

	// documentation for developers
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	port := os.Getenv("PORT")

	fmt.Println("Listening on", port)

	// Start server
	e.Logger.Fatal(e.Start(":" + port))

}
