package models

import (
	"ankasa-be/src/configs"
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	TicketID        uint        `json:"ticket_id" validate:"required"`
	Ticket          Ticket      `gorm:"foreignKey:TicketID" json:"ticket"`
	CustomerID      uint        `json:"customer_id" validate:"required"`
	Customer        Customer    `gorm:"foreignKey:CustomerID" json:"customer"`
	TotalPassengers int         `json:"total_passengers" validate:"required"`
	TotalPrice      float64     `json:"total_price" validate:"required"`
	Passengers      []Passenger `json:"passengers" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ExternalID      string      `json:"external_id"`
	PaymentID       string      `json:"payment_id"`
	PaymentMethod   string      `json:"payment_method" validate:"required"`
	PaymentStatus   string      `json:"payment_status" validate:"required"`
	PaymentURL      string      `json:"payment_url" validate:"required"`
	PaidAt          *time.Time  `json:"paid_at"`
	IsActive        bool        `json:"is_active" validate:"required"`
}

type OrderRequest struct {
	TicketID             uint                  `json:"ticket_id" validate:"required"`
	Passengers           []Passenger           `json:"passengers" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TotalPrice           float64               `json:"total_price" validate:"required"`
	PaymentMethodType    string                `json:"payment_method_type" validate:"required"`
	PaymentMethodCard    *PaymentMethodCard    `json:"payment_method_card"`
	PaymentMethodEWallet *PaymentMethodEWallet `json:"payment_method_ewallet"`
}

type PaymentMethodCard struct {
	CardNumber     string  `json:"card_number"`
	ExpiryMonth    string  `json:"expiry_month"`
	ExpiryYear     string  `json:"expiry_year"`
	CardholderName *string `json:"cardholder_name,omitempty"`
	Cvv            *string `json:"cvv,omitempty"`
}
type PaymentMethodEWallet struct {
	Name         string  `json:"name"`
	MobileNumber *string `json:"mobile_number,omitempty"`
	Cashtag      *string `json:"cashtag,omitempty"`
}

type Passenger struct {
	gorm.Model
	OrderID     uint   `json:"order_id" validate:"required"`
	Order       Order  `gorm:"foreignKey:OrderID" json:"order"`
	SeatID      uint   `json:"seat_id" validate:"required"`
	Seat        Seat   `gorm:"foreignKey:SeatID" json:"seat"`
	Name        string `json:"name" validate:"required,max=50"`
	Category    string `json:"category" validate:"required,max=50"`
	Nationality string `json:"nationality" validate:"required,max=50"`
}

func CreateOrder(order *Order) (uint, error) {
	result := configs.DB.Create(&order)
	return order.ID, result.Error
}

func SelectAllOrders() []*Order {
	var orders []*Order

	configs.DB.Model(&Order{}).
		Preload("Ticket", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Merchant").Preload("Arrival", func(db *gorm.DB) *gorm.DB {
				return db.Preload("City", func(db *gorm.DB) *gorm.DB {
					return db.Preload("Country")
				})
			}).Preload("Departure", func(db *gorm.DB) *gorm.DB {
				return db.Preload("City", func(db *gorm.DB) *gorm.DB {
					return db.Preload("Country")
				})
			})
		}).
		Preload("Customer").
		Find(&orders)

	return orders
}

func SelectOrdersbyCustomerID(customer_id int) []*Order {
	var orders []*Order

	configs.DB.Model(&Order{}).
		Preload("Ticket", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Merchant").Preload("Arrival", func(db *gorm.DB) *gorm.DB {
				return db.Preload("City", func(db *gorm.DB) *gorm.DB {
					return db.Preload("Country")
				})
			}).Preload("Departure", func(db *gorm.DB) *gorm.DB {
				return db.Preload("City", func(db *gorm.DB) *gorm.DB {
					return db.Preload("Country")
				})
			})
		}).
		Preload("Customer").
		Find(&orders, "customer_id", customer_id)

	return orders
}

func SelectOrderbyID(id *int) *Order {
	var order *Order
	configs.DB.Preload("Ticket", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Merchant").Preload("Class").Preload("Arrival", func(db *gorm.DB) *gorm.DB {
			return db.Preload("City", func(db *gorm.DB) *gorm.DB {
				return db.Preload("Country")
			})
		}).Preload("Departure", func(db *gorm.DB) *gorm.DB {
			return db.Preload("City", func(db *gorm.DB) *gorm.DB {
				return db.Preload("Country")
			})
		})
	}).Preload("Customer").Preload("Passengers", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Seat")
	}).First(&order, "id = ?", &id)
	return order
}

func SelectOrderSingle(name string, value interface{}) *Order {
	var order *Order
	configs.DB.First(&order, name, value)
	return order
}

func UpdateOrderById(id int, updatedOrder Order) (int, error) {
	result := configs.DB.Model(&Order{}).Where("id = ?", id).Updates(updatedOrder)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(result.RowsAffected), nil
}

func UpdateOrderSingle(id int, name string, value interface{}) error {
	result := configs.DB.Model(&Order{}).Where("id = ?", id).Update(name, value)
	return result.Error
}

func DeleteOrderById(id int) (int, error) {
	result := configs.DB.Delete(&Order{}, "id = ?", id)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(result.RowsAffected), nil
}
