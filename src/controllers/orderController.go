package controllers

import (
	"ankasa-be/src/helpers"
	"ankasa-be/src/middlewares"
	"ankasa-be/src/models"
	"ankasa-be/src/services"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/xendit/xendit-go/v6/invoice"
)

func GetAllOrders(c *fiber.Ctx) error {
	orders := models.SelectAllOrders()

	if len(orders) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":     "success",
			"statusCode": fiber.StatusOK,
			"message":    "orders unavailable",
		})
	}

	resultOrders := make([]map[string]interface{}, len(orders))
	for i, order := range orders {
		resultOrders[i] = map[string]interface{}{
			"id":                     order.ID,
			"created_at":             order.CreatedAt,
			"updated_at":             order.UpdatedAt,
			"merchant_name":          order.Ticket.Merchant.Name,
			"departure_schedule":     order.Ticket.Departure.Schedule,
			"departure_country_code": order.Ticket.Departure.City.Country.Code,
			"arrival_country_code":   order.Ticket.Arrival.City.Country.Code,
			"payment_method":         order.PaymentMethod,
			"payment_status":         order.PaymentStatus,
			"payment_url":            order.PaymentURL,
			"is_active":              order.IsActive,
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     "success",
		"statusCode": fiber.StatusOK,
		"data":       resultOrders,
	})
}

func CreateOrder(c *fiber.Ctx) error {

	user_id, err := middlewares.JWTAuthorize(c, "customer")
	if err != nil {
		if fiberErr, ok := err.(*fiber.Error); ok {
			return c.Status(fiberErr.Code).JSON(fiber.Map{
				"status":     fiberErr.Message,
				"statusCode": fiberErr.Code,
				"message":    fiberErr.Message,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "Internal Server Error",
			"statusCode": fiber.StatusInternalServerError,
			"message":    err.Error(),
		})
	}

	user := models.SelectUserfromID(int(user_id))
	if user.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    "User not found",
		})
	}

	customer := models.SelectCustomerfromUserID(int(user_id))
	if customer.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    "Customer not found",
		})
	}

	var order models.Order
	if err := c.BodyParser(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Invalid request body",
		})
	}

	if err := models.CreateOrder(&order); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "internal server error",
			"statusCode": fiber.StatusInternalServerError,
			"message":    err.Error(),
		})
	}

	customerDetail := *invoice.NewCustomerObject()
	customerDetail.SetGivenNames(customer.Username)
	customerDetail.SetEmail(user.Email)
	customerDetail.SetMobileNumber(customer.PhoneNumber)

	currency := "IDR"

	externalID := "ankasa-" + helpers.GenerateString(16) + "-" + time.Now().Format("20060102150405")
	createInvoiceRequest := *invoice.NewCreateInvoiceRequest(externalID, float64(order.TotalPrice))
	createInvoiceRequest.PayerEmail = &user.Email
	createInvoiceRequest.Customer = &customerDetail
	createInvoiceRequest.Currency = &currency

	resp, _, _ := services.Client.InvoiceApi.
		CreateInvoice(context.Background()).
		CreateInvoiceRequest(createInvoiceRequest).
		Execute()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":     "created",
		"statusCode": fiber.StatusCreated,
		"message":    "Order has been created successfully",
		"data":       resp,
	})
}
