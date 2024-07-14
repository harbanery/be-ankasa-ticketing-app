package configs

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	url := "postgres://default:Y1xpNs9bQwHT@ep-plain-firefly-a186xywd-pooler.ap-southeast-1.aws.neon.tech:5432/verceldb?sslmode=require"
	var err error
	DB, err = gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
}
