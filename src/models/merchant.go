package models

import (
	"ankasa-be/src/configs"

	"gorm.io/gorm"
)

type Merchant struct {
	gorm.Model
	UserID      uint    `json:"user_id" validate:"required"`
	User        User    `gorm:"foreignKey:UserID" validate:"-"`
	Name        string  `json:"name" validate:"required,max=50"`
	Image       string  `json:"image" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Classes     []Class `json:"classes"`
}

type Class struct {
	gorm.Model
	MerchantID     uint   `json:"merchant_id" validate:"required"`
	Name           string `json:"name" validate:"required,max=50"`
	Quantity       uint   `json:"quantity" validate:"required"`
	Price          uint   `json:"price" validate:"required"`
	IsRefund       bool   `json:"is_refund"`
	IsReschedule   bool   `json:"is_reschedule"`
	IsLuggage      bool   `json:"is_luggage"`
	IsInflightMeal bool   `json:"is_inflight_meal"`
	IsWifi         bool   `json:"is_wifi"`
}

func CreateMerchant(merchant *Merchant) (uint, error) {
	result := configs.DB.Create(&merchant)
	return merchant.ID, result.Error
}

func CreateClass(class *Class) error {
	result := configs.DB.Create(&class)
	return result.Error
}

func SelectAllMerchants() []*Merchant {
	var merchants []*Merchant
	configs.DB.Preload("Classes").Find(&merchants)
	return merchants
}
