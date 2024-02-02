package database

import (
	"log"
	"os"

	"github.com/sixfwa/fiber-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbInstance struct {
	Db *gorm.DB
}

var Database DbInstance

func ConnectDb() {
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database! \n", err.Error())
		os.Exit(3)
	}

	log.Println("Connected to the database successfully")
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("Running Migrations")
	// Add migrations
	db.AutoMigrate(&models.Year{}, &models.Item{}, &models.Test{}, &models.Question{})

	Database = DbInstance{Db: db}
}
