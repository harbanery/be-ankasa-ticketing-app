package models

import (
	"ankasa-be/src/configs"
	"time"

	"gorm.io/gorm"
)

type Ticket struct {
	gorm.Model
	Stock      int        `json:"stock" validate:"required"`
	Price      float64    `json:"price" validate:"required"`
	ClassID    uint       `json:"class_id" validate:"required"`
	Class      Class      `gorm:"foreignKey:ClassID" json:"class"`
	Gate       string     `json:"gate" validate:"required"`
	Categories []Category `json:"categories"`
	Arrival    Arrival    `json:"arrival"`
	Departure  Departure  `json:"departure"`
	// Transit        Transit    `json:"transit"`
}

type Arrival struct {
	gorm.Model
	Schedule *time.Time `json:"schedule" validate:"required"`
	CityID   uint       `json:"city_id" validate:"required"`
	TicketID uint       `json:"ticket_id" validate:"required"`
}

type Departure struct {
	gorm.Model
	Schedule *time.Time `json:"schedule" validate:"required"`
	CityID   uint       `json:"city_id" validate:"required"`
	TicketID uint       `json:"ticket_id" validate:"required"`
}

type Transit struct {
	gorm.Model
	Name     string `json:"name" validate:"required"`
	CityID   uint   `json:"city_id" validate:"required"`
	TicketID uint   `json:"ticket_id" validate:"required"`
}

func SelectAllTickets() ([]Ticket, error) {
	var tickets []Ticket

	err := configs.DB.Model(&Ticket{}).Preload("Class").Preload("Categories").Preload("Arrival").Preload("Departure").Find(&tickets).Error

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
