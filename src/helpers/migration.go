package helpers

import (
	"ankasa-be/src/configs"
	"ankasa-be/src/models"
	"log"
)

func Migration() {
	err := configs.DB.AutoMigrate(&models.User{})

	if err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}
}
