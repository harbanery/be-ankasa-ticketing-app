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

func SelectCustomerfromID(id int) *Customer {
	var customer Customer
	configs.DB.Preload("User").First(&customer, "id = ?", id)
	return &customer
}

func SelectCustomerfromUserID(user_id int) *Customer {
	var customer Customer
	configs.DB.Preload("User").First(&customer, "user_id = ?", user_id)
	return &customer
}

func UpdateCustomer(id int, updatedCustomer *Customer) error {
	result := configs.DB.Model(&Customer{}).Where("id = ?", id).Updates(updatedCustomer)
	return result.Error
}

func UpdateCustomerSingle(id int, name, value string) error {
	result := configs.DB.Model(&Customer{}).Where("id = ?", id).Update(name, value)
	return result.Error
}
