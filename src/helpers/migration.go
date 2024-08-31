package helpers

import (
	"ankasa-be/src/configs"
	"ankasa-be/src/models"
	"log"
)

func Migration() {
	err := configs.DB.AutoMigrate(
		&models.User{},
		&models.Merchant{},
		&models.Customer{},
		&models.UserVerification{},
		&models.Ticket{},
		&models.Country{},
		&models.City{},
		&models.Category{},
		&models.Chat{},
		&models.ChatUser{},
		&models.Message{},
	)

	if err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}
}
