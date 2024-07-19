package models

import (
	"ankasa-be/src/configs"

	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	UserID      uint   `json:"user_id" validate:"required"`
	User        User   `gorm:"foreignKey:UserID" validate:"-"`
	Username    string `json:"username" validate:"required,max=50"`
	PhoneNumber string `json:"phone_number" validate:"required,numeric,max=15"`
	City        string `json:"city"  validate:"required"`
	Image       string `json:"image" validate:"required"`
	Address     string `json:"address" validate:"required"`
	PostalCode  string `json:"postal_code" validate:"required,numeric,max=5"`
}

func CreateCustomer(customer *Customer) error {
	result := configs.DB.Create(&customer)
	return result.Error
}
