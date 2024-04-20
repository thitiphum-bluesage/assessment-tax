package infrastructure

import (
	"log"

	"github.com/thitiphum-bluesage/assessment-tax/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitializeDatabase() *gorm.DB {

	cfg := config.GetConfig()
	DATABASE_URI := cfg.DatabaseURL

	db, err := gorm.Open(postgres.Open(DATABASE_URI), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Successfully connected to database.")
	return db

}
