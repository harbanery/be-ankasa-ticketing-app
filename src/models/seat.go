package models

import (
	"ankasa-be/src/configs"

	"gorm.io/gorm"
)

type Seat struct {
	gorm.Model
	Code      string `json:"code" validate:"required"`
	IsBooking bool   `json:"is_booking"`
	TicketID  uint   `json:"ticket_id" validate:"required"`
}

func CreateSeat(seat *Seat) error {
	err := configs.DB.Create(&seat).Error
	return err
}

func UpdateSeatIsBooking(id int, value bool) error {
	result := configs.DB.Model(&Seat{}).Where("id = ?", id).Update("is_booking", value)
	return result.Error
}
