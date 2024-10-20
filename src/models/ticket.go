package models

import (
	"ankasa-be/src/configs"
	"log"
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

// func SelectAllTicketWithFilter(sort string, limit, offset int, filter map[string]interface{}, condition string) []*Ticket {
// 	var ticket []*Ticket

// 	query := configs.DB.Preload("Merchant").Preload("Class").
// 		Preload("Arrival", func(db *gorm.DB) *gorm.DB {
// 			return db.Preload("City", func(db *gorm.DB) *gorm.DB {
// 				return db.Preload("Country")
// 			})
// 		}).Preload("Departure", func(db *gorm.DB) *gorm.DB {
// 		return db.Preload("City", func(db *gorm.DB) *gorm.DB {
// 			return db.Preload("Country")
// 		})
// 	}).Group("ticket.id").Order(sort).Limit(limit).Offset(offset)

// 	if merchant, ok := filter["merchatValues"].([]string); ok && len(merchant) > 0 {
// 		query = query.Joins("INNER JOIN merchant ON merchant.ticket_id = ticket.id AND merchant.value IN ?", filter["merchantValues"])
// 	}

// 	if arrival, ok := filter["arrivalValues"].([]string); ok && len(arrival) > 0 {
// 		query = query.Joins("INNER JOIN arrival ON arrival.ticket_id = ticket.id AND arrival.value IN ?", filter["arrivalValues"])
// 	}

// 	if departure, ok := filter["departureValues"].([]string); ok && len(departure) > 0 {
// 		query = query.Joins("INNER JOIN departure ON departure.ticket_id = ticket.id AND departure.value IN ?", filter["departureValues"])
// 	}

// 	if class, ok := filter["classValues"].([]string); ok && len(class) > 0 {
// 		query = query.Joins("INNER JOIN class ON class.ticket_id = ticket.id AND class.value IN ?", filter["classValues"])
// 	}

// 	if condition != "" {
// 		query = query.Where("ticket.condition = ?", condition)
// 	}

// 	query.Find(&ticket)

// 	return ticket
// }
func SelectAllTicketWithFilter(sort string, limit, offset int, filter map[string]interface{}, condition string) []*Ticket {
    var ticket []*Ticket

    // Secara eksplisit tentukan tabel tickets
    query := configs.DB.Table("tickets").Preload("Merchant").Preload("Class").
        Preload("Arrival", func(db *gorm.DB) *gorm.DB {
            return db.Preload("City", func(db *gorm.DB) *gorm.DB {
                return db.Preload("Country")
            })
        }).Preload("Departure", func(db *gorm.DB) *gorm.DB {
        return db.Preload("City", func(db *gorm.DB) *gorm.DB {
            return db.Preload("Country")
        })
    }).Group("tickets.id").Order(sort).Limit(limit).Offset(offset)

    // Filter merchant
    if merchant, ok := filter["merchantValues"].([]string); ok && len(merchant) > 0 {
        query = query.Joins("JOIN merchants ON merchants.id = tickets.merchant_id AND merchants.name IN ?", merchant)
    }

    // Filter arrival
    if arrival, ok := filter["arrivalValues"].([]time.Time); ok && len(arrival) > 0 {
        query = query.Joins("JOIN arrivals ON arrivals.ticket_id = tickets.id AND arrivals.schedule IN ?", arrival)
    }

    // Filter departure
    if departure, ok := filter["departureValues"].([]time.Time); ok && len(departure) > 0 {
        query = query.Joins("JOIN departures ON departures.ticket_id = tickets.id AND departures.schedule IN ?", departure)
    }

    // Filter class
    if class, ok := filter["classValues"].([]string); ok && len(class) > 0 {
        query = query.Joins("JOIN classes ON classes.id = tickets.class_id AND classes.name IN ?", class)
    }

    // Filter condition
    if condition != "" {
        query = query.Where("tickets.condition = ?", condition)
    }

    // Eksekusi query
    result := query.Find(&ticket)
    if result.Error != nil {
        log.Printf("Error fetching tickets: %v", result.Error)
    }

    return ticket
}


