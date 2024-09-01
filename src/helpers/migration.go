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
		&models.Class{},
		&models.Ticket{},
		&models.Arrival{},
		&models.Departure{},
		&models.Country{},
		&models.City{},
		&models.Seat{},
		&models.Chat{},
		&models.ChatUser{},
		&models.Message{},
		&models.Wallet{},
	)

	if err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}
}
