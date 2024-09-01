package models

import (
	"ankasa-be/src/configs"

	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name     string `json:"name" validate:"required"`
	Quantity int    `json:"quantity" validate:"required"`
	TicketID uint   `json:"ticket_id" validate:"required"`
	Seats    []Seat `json:"seats"`
}

func SelectAllCategories() ([]Category, error) {
	var categories []Category

	err := configs.DB.Find(&categories).Error

	if err != nil {
		return nil, err
	}

	return categories, nil
}

func SelectCategoryById(id int) (Category, error) {
	var category Category
	if err := configs.DB.First(&category, "id = ?", id).Error; err != nil {
		return Category{}, err
	}

	return category, nil
}

func CreateCategory(category *Category) (uint, error) {
	err := configs.DB.Create(&category).Error
	return category.ID, err
}

func UpdateCategoryById(id int, updatedCategory Category) (int, error) {
	result := configs.DB.Model(&Category{}).Where("id = ?", id).Updates(updatedCategory)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(result.RowsAffected), nil
}

func DeleteCategoryById(id int) (int, error) {
	result := configs.DB.Delete(&Category{}, "id = ?", id)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(result.RowsAffected), nil
}
