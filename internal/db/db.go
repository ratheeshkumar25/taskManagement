package db

import (
	"log"

	"github.com/ratheeshkumar25/task-mgt/config"
	"github.com/ratheeshkumar25/task-mgt/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB(config *config.Config) *gorm.DB {
	//log.Printf("Database URL from config: %s", config.Database_url)

	// if config.Database_url == "" {
	// 	log.Fatal("Database URL is not set in the configuration")
	// }
	// // Log the connection details
	// log.Printf("Connecting to DB with URL: %s", config.Database_url)

	DB, err := gorm.Open(postgres.Open(config.Database_url), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Printf("Error while connecting to the database: %v", err)
		return nil
	}

	// Migrate the schema
	if err := DB.AutoMigrate(&models.Users{}, &models.Task{}); err != nil {
		log.Printf("Error while migrating: %v", err)
		return nil
	}

	return DB
}
