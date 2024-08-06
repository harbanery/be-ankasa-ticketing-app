package models

import (
	"ankasa-be/src/configs"

	"gorm.io/gorm"
)

type Ticket struct {
	gorm.Model
	Stock          uint       `json:"stock" validate:"required"`
	Price          uint       `json:"price" validate:"required"`
	Class          string     `json:"class" validate:"required"`
	Gate           string     `json:"gate" validate:"required"`
	IsRefund       bool       `json:"is_refund"`
	IsReschedule   bool       `json:"is_reschedule"`
	IsLuggage      bool       `json:"is_luggage"`
	IsInflightMeal bool       `json:"is_inflight_meal"`
	IsWifi         bool       `json:"is_wifi"`
	Categories     []Category `json:"categories"`
}

func SelectAllTickets() ([]Ticket, error) {
	var tickets []Ticket

	err := configs.DB.Model(&Ticket{}).Preload("Categories").Find(&tickets).Error

	if err != nil {
		return nil, err
	}

	return tickets, nil
}

func SelectTicketById(id int) (Ticket, error) {
	var ticket Ticket
	if err := configs.DB.Model(&Ticket{}).Preload("Categories").First(&ticket, "id = ?", id).Error; err != nil {
		return Ticket{}, err
	}

	return ticket, nil
}

func CreateTicket(ticket *Ticket) error {
	err := configs.DB.Create(&ticket).Error
	return err
}

func UpdateTicketById(id int, updatedTicket Ticket) (int, error) {
	result := configs.DB.Model(&Ticket{}).Where("id = ?", id).Updates(updatedTicket)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(result.RowsAffected), nil
}

func DeleteTicketById(id int) (int, error) {
	result := configs.DB.Delete(&Ticket{}, "id = ?", id)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(result.RowsAffected), nil
}
