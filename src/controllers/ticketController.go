package controllers

import (
	"ankasa-be/src/models"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetAllTickets(c *fiber.Ctx) error {
	tickets := models.SelectAllTickets()

	if len(tickets) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    "tickets not found",
		})
	}

	resultTickets := make([]map[string]interface{}, len(tickets))
	for i, ticket := range tickets {
		resultTickets[i] = map[string]interface{}{
			"id":                     ticket.ID,
			"created_at":             ticket.CreatedAt,
			"updated_at":             ticket.UpdatedAt,
			"merchant_name":          ticket.Merchant.Name,
			"merchant_image":         ticket.Merchant.Image,
			"departure_schedule":     ticket.Departure.Schedule,
			"departure_country_code": ticket.Departure.City.Country.Code,
			"arrival_schedule":       ticket.Arrival.Schedule,
			"arrival_country_code":   ticket.Arrival.City.Country.Code,
			"price":                  ticket.Price,
			"is_refund":              ticket.Class.IsRefund,
			"is_reschedule":          ticket.Class.IsReschedule,
			"is_luggage":             ticket.Class.IsLuggage,
			"is_inflight_meal":       ticket.Class.IsInflightMeal,
			"is_wifi":                ticket.Class.IsWifi,
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     "success",
		"statusCode": fiber.StatusOK,
		"data":       resultTickets,
	})
}

func GetTicketById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Invalid ID format",
		})
	}

	ticket, err := models.SelectTicketById(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    "Ticket is not found",
		})
	}

	resultSeats := make([]map[string]interface{}, len(ticket.Seats))
	for i, seat := range ticket.Seats {
		resultSeats[i] = map[string]interface{}{
			"id":         seat.ID,
			"code":       seat.Code,
			"is_booking": seat.IsBooking,
		}
	}

	resultTickets := map[string]interface{}{
		"id":                     ticket.ID,
		"created_at":             ticket.CreatedAt,
		"updated_at":             ticket.UpdatedAt,
		"merchant_name":          ticket.Merchant.Name,
		"merchant_image":         ticket.Merchant.Image,
		"departure_schedule":     ticket.Departure.Schedule,
		"departure_city":         ticket.Departure.City.Name,
		"departure_country_code": ticket.Departure.City.Country.Code,
		"arrival_schedule":       ticket.Arrival.Schedule,
		"arrival_city":           ticket.Arrival.City.Name,
		"arrival_country_code":   ticket.Arrival.City.Country.Code,
		"price":                  ticket.Price,
		"class":                  ticket.Class.Name,
		"is_refund":              ticket.Class.IsRefund,
		"is_reschedule":          ticket.Class.IsReschedule,
		"is_luggage":             ticket.Class.IsLuggage,
		"is_inflight_meal":       ticket.Class.IsInflightMeal,
		"is_wifi":                ticket.Class.IsWifi,
		"row_seats":              ticket.Class.RowSeats,
		"seats":                  resultSeats,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     "success",
		"statusCode": fiber.StatusOK,
		"data":       resultTickets,
	})
}

func CreateTicket(c *fiber.Ctx) error {
	var ticket models.Ticket
	if err := c.BodyParser(&ticket); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Invalid request body",
		})
	}

	if _, err := models.CreateTicket(&ticket); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "internal server error",
			"statusCode": fiber.StatusInternalServerError,
			"message":    err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":     "created",
		"statusCode": fiber.StatusCreated,
		"message":    "Ticket has been created successfully",
		"data":       fiber.Map{"id": ticket.ID},
	})
}

func UpdateTicket(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Invalid ID format",
		})
	}

	var updatedTicket models.Ticket

	if err := c.BodyParser(&updatedTicket); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Invalid request body",
			"data":       updatedTicket,
		})
	}

	rowsAffected, err := models.UpdateTicketById(id, updatedTicket)
	if err != nil || rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    fmt.Sprintf("Ticket with ID %d is not found", id),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     "success",
		"statusCode": 200,
		"message":    fmt.Sprintf("Ticket with ID %d has been updated successfully", id),
	})
}

func DeleteTicket(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Invalid ID format",
		})
	}
	rowsAffected, err := models.DeleteTicketById(id)
	if err != nil || rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    fmt.Sprintf("Ticket with ID %d is not found", id),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     "success",
		"statusCode": 200,
		"message":    fmt.Sprintf("Ticket with ID %d has been deleted successfully", id),
	})
}
