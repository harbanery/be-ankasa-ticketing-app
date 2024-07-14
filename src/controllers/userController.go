package controllers

import (
	"ankasa-be/src/models"

	"github.com/gofiber/fiber/v2"
)

func GetUsers(c *fiber.Ctx) error {
	users := models.SelectUsers()
	if len(users) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    "Users not found.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     "success",
		"statusCode": 200,
		"data":       users,
	})
}

func RegisterUser(c *fiber.Ctx) error {
	var register models.User
	if err := c.BodyParser(&register); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Invalid request body",
		})
	}

	err := models.CreateUser(&register)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":     "success",
		"statusCode": 200,
		"message":    "User created successfully.",
	})
}
