package db

import (
	"AiCompServer/app/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

func InitDB() {
	dbhost := os.Getenv("DBHOST")
	// dbname := os.Getenv("DBNAME")
	// dbuser := os.Getenv("DBUSER")
	// dbpass := os.Getenv("DBPASSWORD")
	for {
		db, err := gorm.Open("postgres", "host="+dbhost+" port=5432 user=gorm dbname=gorm sslmode=disable password=yatuhashi-api")
		if err != nil {
			log.Println("Failed to connect to database: %v\n", err)
		} else {
			db.DB()
			db.AutoMigrate(&models.User{})
			db.AutoMigrate(&models.Challenge{})
			db.AutoMigrate(&models.Answer{})
			db.AutoMigrate(&models.PreAnswer{})
			db.AutoMigrate(&models.Notification{})
			db.AutoMigrate(&models.Setting{})
			DB = db
			break
		}
		time.Sleep(1 * time.Second)
	}
}
