package database

import (
	"github.com/romitou/disneystats/database/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

var database *gorm.DB

func ConnectDatabase() {
	mysqlDsn := os.Getenv("MYSQL_DSN")
	db, err := gorm.Open(mysql.Open(mysqlDsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&models.Attraction{}, &models.AttractionWaitTime{})
	if err != nil {
		log.Println("an error occurred while migrating the database: ", err)
	}

	database = db
}

func GetDatabase() *gorm.DB {
	return database
}
