package infrastructure

import (
	"log"

	"github.com/thitiphum-bluesage/assessment-tax/config"
	"github.com/thitiphum-bluesage/assessment-tax/domains"
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

	if err := db.AutoMigrate(&domains.TaxDeductionConfig{}); err != nil {
        log.Fatalf("Failed to migrate database: %v", err)
    }

	if err := ensureDefaultConfigExists(db); err != nil {
		log.Fatalf("Failed to initialize default configuration: %v", err)
	}

	log.Println("Successfully connected to database.")
	return db
}

func ensureDefaultConfigExists(db *gorm.DB) error {
	var count int64
	if err := db.Model(&domains.TaxDeductionConfig{}).Count(&count).Error; err != nil {
		return err  
	}
	if count == 0 {
		defaultConfig := domains.TaxDeductionConfig{
			ConfigName: "MainConfig",
			PersonalDeduction: 60000,
			KReceiptDeductionMax: 50000,
			DonationDeductionMax: 100000,
		}
		if err := db.Create(&defaultConfig).Error; err != nil {
			return err 
		}
		log.Println("Main tax deduction configuration has been initialized.")
	}
	return nil  
}
