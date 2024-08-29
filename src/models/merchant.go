package models

import (
	"ankasa-be/src/configs"

	"gorm.io/gorm"
)

type Merchant struct {
	gorm.Model
	UserID      uint   `json:"user_id" validate:"required"`
	User        User   `gorm:"foreignKey:UserID" validate:"-"`
	Name        string `json:"name" validate:"required,max=50"`
	Image       string `json:"image" validate:"required"`
	Description string `json:"description" validate:"required"`
}

func CreateMerchant(merchant *Merchant) error {
	result := configs.DB.Create(&merchant)
	return result.Error
}
