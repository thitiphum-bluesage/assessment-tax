package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/thitiphum-bluesage/assessment-tax/applications/services"
	"github.com/thitiphum-bluesage/assessment-tax/config"
	"github.com/thitiphum-bluesage/assessment-tax/infrastructure"
	"github.com/thitiphum-bluesage/assessment-tax/infrastructure/repository"
	"github.com/thitiphum-bluesage/assessment-tax/interfaces/endpoints"
	"github.com/thitiphum-bluesage/assessment-tax/interfaces/endpoints/controllers"
)

func main() {

	// Load configuration
	cfg := config.GetConfig()

	db := infrastructure.InitializeDatabase()

	fmt.Println(db)

	e := echo.New()

	// Repository layer
	taxRepo := repository.NewTaxDeductionConfigRepository(db)

	// Service layer
	adminService := services.NewAdminService(taxRepo)
	taxService := services.NewTaxService(taxRepo)

	// Controller layer
	adminController := controllers.NewAdminController(adminService)
	taxController := controllers.NewTaxController(taxService)

	// Setup the router with routes
	endpoints.NewRouter(e, taxController, adminController, cfg)

	port := cfg.Port
	if port == "" {
		log.Fatal("PORT environment variable not set")
	}

	// Set up Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("shutting down the server: %v", err)
		}
	}()

	// Block until a signal is received
	<-quit

	handleGracefulShutdown(e)
}

func handleGracefulShutdown(e *echo.Echo) {
	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt to shut down the server gracefully
	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to gracefully shut down the server: %v", err)
	}

	log.Println("Server shut down gracefully")
}
