package models

import (
	"ankasa-be/src/configs"

	"gorm.io/gorm"
)

type Wallet struct {
	gorm.Model
	Saldo      float64  `json:"saldo" validate:"min=0"`
	CustomerID int      `json:"customer_id" validate:"required"`
	Customer   Customer `gorm:"foreignKey:CustomerID"`
}

type APIWallet struct {
	gorm.Model
	Saldo      float64 `json:"saldo" validate:"min=0"`
	CustomerID int     `json:"customer_id"`
}

func GetAllWallet() ([]*Wallet, error) {
	var wallet []*Wallet
	results := configs.DB.Preload("Customer").Find(&wallet)
	if results.Error != nil {
		return nil, results.Error
	}
	return wallet, nil
}

func GetWalletById(id int) *Wallet {
	var wallet *Wallet
	configs.DB.Preload("Customer").First(&wallet, "id = ?", id)
	return wallet
}

func GetWalletByCustomerId(customer_id int) *Wallet {
	var wallet *Wallet
	configs.DB.Preload("Customer").First(&wallet, "customer_id = ?", customer_id)
	return wallet
}

func CreateWallet(wallet *Wallet) error {
	result := configs.DB.Create(&wallet)
	return result.Error
}

func UpdateWallet(id int, newWalet *Wallet) error {
	result := configs.DB.Model(&Wallet{}).Where("id = ?", id).Updates(&newWalet)
	return result.Error
}

func UpdateWalletSaldo(id int, saldo float64) error {
	result := configs.DB.Model(&Ticket{}).Where("id = ?", id).Update("saldo", saldo)
	return result.Error
}

func DeleteWallet(id int) error {
	result := configs.DB.Where("id = ? ", id).Delete(&Wallet{})
	return result.Error
}
