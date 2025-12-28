package config

import (
	"log"
	"os"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func GetDBUrl() string {
	_ = godotenv.Load()

	url := os.Getenv("DATABASE_URL")
	if url == ""{
		log.Fatal("DATABASE URL NOT SET")
	}

	return url
}

func ConnectDB() {
	var err error
	dsn := GetDBUrl()
	
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	
	log.Println("Database connected successfully!")
}
