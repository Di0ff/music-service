package db

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"music-service/internal/model"
	"os"
)

var DB *gorm.DB

func Init() {
	var errDB error

	errConfig := godotenv.Load("config.env")
	if errConfig != nil {
		log.Fatal("Error loading config.env")
	}

	url := os.Getenv("DATABASE_URL")
	if url == "" {
		log.Fatal("DATABASE_URL environment variable not set")
		return
	}

	DB, errDB = gorm.Open(postgres.Open(url), &gorm.Config{})
	if errDB != nil {
		log.Fatal("Error connecting to the database: ", errDB)
	}

	errDB = DB.AutoMigrate(&model.Song{})
	if errDB != nil {
		log.Fatal("Error during database migration:", errDB)
	}
}
