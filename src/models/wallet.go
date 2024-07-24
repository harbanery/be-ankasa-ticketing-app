package models

import (
	"ankasa-be/src/configs"

	"gorm.io/gorm"
)

type Wallet struct {
	gorm.Model
	Saldo      float64  `json:"saldo" validate:"min=0"`
	CustomerID uint     `json:"customer_id"`
	Customer   Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
}

func GetAllWallet() ([]*Wallet, error) {
	var wallet []*Wallet
	results := configs.DB.Find(&wallet)
	if results.Error != nil {
		return nil, results.Error
	}
	return wallet, nil
}

func GetWalletById(id uint) *Wallet {
	var wallet *Wallet
	configs.DB.First(&wallet, "id = ?", id)
	return wallet
}

func CreateWallet(wallet *Wallet) error {
	result := configs.DB.Create(&wallet)
	return result.Error
}

func UpdateWallet(id uint, newWalet *Wallet) error {
	result := configs.DB.Model(&Wallet{}).Where("id = ?", id).Updates(&newWalet)
	return result.Error
}

func DeleteWallet(id uint) error {
	result := configs.DB.Where("id = ? ", id).Delete(&Wallet{})
	return result.Error
}
