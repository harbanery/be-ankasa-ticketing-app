package controllers

import (
	"ankasa-be/src/helpers"
	"ankasa-be/src/middlewares"
	"ankasa-be/src/models"

	"github.com/gofiber/fiber/v2"
)
// Buat testing aja
func GetCustomers(c *fiber.Ctx) error {
	customer := models.SelectCustomers()
	if len(customer) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    "customer not found.",
		})
	}
	return c.Status(200).JSON(customer)
}
// 

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

func UpdateCustomerProfile(c *fiber.Ctx) error {
	var profileData models.CustomerProfile

	id, err := middlewares.JWTAuthorize(c, "customer")
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

	customer := models.SelectCustomerfromUserID(int(id))
	if customer.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "Not Found",
			"statusCode": 404,
			"message":    "Customer not found",
		})
	}

	if err := c.BodyParser(&profileData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Invalid Request Body",
		})
	}

	user := middlewares.XSSMiddleware(&profileData).(*models.CustomerProfile)
	if errors := helpers.StructValidation(user); len(errors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":     "unprocessable entity",
			"statusCode": "422",
			"message":    "Validation failed",
			"error":      errors,
		})
	}

	if customer.User.Email != user.Email {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "bad request",
			"statusCode": 400,
			"message": "Email already exists",
		})
	}

	updatedCustomer := models.Customer{
		Username: user.Username,
		PhoneNumber: user.PhoneNumber,
		City: user.City,
		Image: user.Image,
		Address: user.Address,
		PostalCode: user.PostalCode,
	}

	if err := models.UpdateUserSingle(int(id), "email", user.Email); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "server error",
			"statusCode": 500,
			"message": "Failed to update user",
		})
	}

	if err := models.UpdateCustomer(int(customer.ID), &updatedCustomer); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "server error",
			"statusCode": 500,
			"message": "Failed to update customer",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"statusCode": 200,
		"message": "Profile updated successfully",
	})

}
