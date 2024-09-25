package models

import (
	"ankasa-be/src/configs"
	"time"

	"gorm.io/gorm"
)

type Ticket struct {
	gorm.Model
	Stock      int       `json:"stock" validate:"required"`
	Price      float64   `json:"price" validate:"required"`
	MerchantID uint      `json:"merchant_id" validate:"required"`
	Merchant   Merchant  `gorm:"foreignKey:MerchantID" json:"merchant"`
	ClassID    uint      `json:"class_id" validate:"required"`
	Class      Class     `gorm:"foreignKey:ClassID" json:"class"`
	Gate       string    `json:"gate" validate:"required"`
	Seats      []Seat    `json:"seats"`
	Arrival    Arrival   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"arrival"`
	Departure  Departure `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"departure"`
	// Transit        Transit    `json:"transit"`
}

type Arrival struct {
	gorm.Model
	Schedule *time.Time `json:"schedule" validate:"required"`
	CityID   uint       `json:"city_id" validate:"required"`
	City     City       `gorm:"foreignKey:CityID" json:"city"`
	TicketID uint       `json:"ticket_id" validate:"required"`
}

type Departure struct {
	gorm.Model
	Schedule *time.Time `json:"schedule" validate:"required"`
	CityID   uint       `json:"city_id" validate:"required"`
	City     City       `gorm:"foreignKey:CityID" json:"city"`
	TicketID uint       `json:"ticket_id" validate:"required"`
}

type Transit struct {
	gorm.Model
	Name     string `json:"name" validate:"required"`
	CityID   uint   `json:"city_id" validate:"required"`
	City     City   `gorm:"foreignKey:CityID" json:"city"`
	TicketID uint   `json:"ticket_id" validate:"required"`
}

func SelectAllTickets() []*Ticket {
	var tickets []*Ticket

	configs.DB.Model(&Ticket{}).
		Preload("Merchant").Preload("Class").
		Preload("Arrival", func(db *gorm.DB) *gorm.DB {
			return db.Preload("City", func(db *gorm.DB) *gorm.DB {
				return db.Preload("Country")
			})
		}).Preload("Departure", func(db *gorm.DB) *gorm.DB {
		return db.Preload("City", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Country")
		})
	}).Order("created_at DESC").Find(&tickets)

	return tickets
}

func SelectTicketById(id int) (Ticket, error) {
	var ticket Ticket
	err := configs.DB.Model(&Ticket{}).
		Preload("Merchant").Preload("Class").Preload("Seats", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at ASC")
	}).Preload("Arrival", func(db *gorm.DB) *gorm.DB {
		return db.Preload("City", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Country")
		})
	}).Preload("Departure", func(db *gorm.DB) *gorm.DB {
		return db.Preload("City", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Country")
		})
	}).First(&ticket, "id = ?", id).Error

	if err != nil {
		return Ticket{}, err
	}

	return ticket, nil
}

func CreateTicket(ticket *Ticket) (uint, error) {
	err := configs.DB.Create(&ticket).Error
	return ticket.ID, err
}

func UpdateTicketById(id int, updatedTicket Ticket) (int, error) {
	result := configs.DB.Model(&Ticket{}).Where("id = ?", id).Updates(updatedTicket)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(result.RowsAffected), nil
}

func UpdateTicketSingle(id int, name string, value interface{}) error {
	result := configs.DB.Model(&Ticket{}).Where("id = ?", id).Update(name, value)
	return result.Error
}

func DeleteTicketById(id int) (int, error) {
	result := configs.DB.Delete(&Ticket{}, "id = ?", id)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(result.RowsAffected), nil
}
