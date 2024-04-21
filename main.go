package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/thitiphum-bluesage/assessment-tax/config"
	"github.com/thitiphum-bluesage/assessment-tax/infrastructure"
)

func main() {

	// Load configuration
	cfg := config.GetConfig()

	db := infrastructure.InitializeDatabase()

	fmt.Println(db)

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})

	port := cfg.Port
	if port == "" {
		log.Fatal("PORT environment variable not set")
	}
	e.Logger.Fatal(e.Start(":" + port))
}
