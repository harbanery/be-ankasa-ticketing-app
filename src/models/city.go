package models

import (
	"ankasa-be/src/configs"

	"gorm.io/gorm"
)

type City struct {
	gorm.Model
	Name      string  `json:"name" validate:"required"`
	Image     string  `json:"image" validate:"required"`
	CountryID uint    `json:"country_id" validate:"required"`
	Country   Country `gorm:"foreignKey:CountryID" json:"country"`
}

func SelectAllCities() ([]City, error) {
	var cities []City

	err := configs.DB.Preload("Country").Find(&cities).Error

	if err != nil {
		return nil, err
	}

	return cities, nil
}

func SelectCityById(id int) (City, error) {
	var city City
	if err := configs.DB.Preload("Country").First(&city, "id = ?", id).Error; err != nil {
		return City{}, err
	}

	return city, nil
}

func CreateCity(city *City) error {
	err := configs.DB.Create(&city).Error
	return err
}

func UpdateCityById(id int, updatedCity City) (int, error) {
	result := configs.DB.Model(&City{}).Where("id = ?", id).Updates(updatedCity)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(result.RowsAffected), nil
}

func DeleteCityById(id int) (int, error) {
	result := configs.DB.Delete(&City{}, "id = ?", id)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(result.RowsAffected), nil
}
