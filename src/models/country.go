package models

import (
	"ankasa-be/src/configs"

	"gorm.io/gorm"
)

type Country struct {
	gorm.Model
	Name   string `json:"name" validate:"required"`
	Code   string `json:"code" validate:"required"`
	Cities []City `json:"cities"`
}

func SelectAllCountries() ([]Country, error) {
	var countries []Country

	err := configs.DB.Model(&Country{}).Preload("Cities").Find(&countries).Error

	if err != nil {
		return nil, err
	}

	return countries, nil
}

func SelectCountryById(id int) (Country, error) {
	var country Country
	if err := configs.DB.Model(&Country{}).Preload("Cities").First(&country, "id = ?", id).Error; err != nil {
		return Country{}, err
	}

	return country, nil
}

func CreateCountry(country *Country) error {
	err := configs.DB.Create(&country).Error
	return err
}

func UpdateCountryById(id int, updatedCountry Country) (int, error) {
	result := configs.DB.Model(&Country{}).Where("id = ?", id).Updates(updatedCountry)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(result.RowsAffected), nil
}

func DeleteCountryById(id int) (int, error) {
	result := configs.DB.Delete(&Country{}, "id = ?", id)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(result.RowsAffected), nil
}
