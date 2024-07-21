package controllers

import (
	"ankasa-be/src/middlewares"
	"ankasa-be/src/models"

	"github.com/gofiber/fiber/v2"
)

func GetCustomerProfile(c *fiber.Ctx) error {
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

	customer := models.SelectCustomerfromUserID(int(user_id))
	if customer.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    "Customer not found",
		})
	}

	resultCustomer := map[string]interface{}{
		"id":           customer.ID,
		"created_at":   customer.CreatedAt,
		"updated_at":   customer.UpdatedAt,
		"username":     customer.Username,
		"user_id":      customer.User.ID,
		"email":        customer.User.Email,
		"image":        customer.Image,
		"phone_number": customer.PhoneNumber,
		"city":         customer.City,
		"address":      customer.Address,
		"postal_code":  customer.PostalCode,
		"role":         customer.User.Role,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     "success",
		"statusCode": 200,
		"data":       resultCustomer,
	})
}
