package controllers

import (
	"ankasa-be/src/models"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetAllCities(c *fiber.Ctx) error {
	cities, err := models.SelectAllCities()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "internal server error",
			"statusCode": 500,
			"message":    err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     "success",
		"statusCode": fiber.StatusOK,
		"data":       cities,
	})
}

func GetCityById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Invalid ID format",
		})
	}

	city, err := models.SelectCityById(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    "City is not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     "success",
		"statusCode": fiber.StatusOK,
		"data":       city,
	})
}

func CreateCity(c *fiber.Ctx) error {
	var city models.City
	if err := c.BodyParser(&city); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Invalid request body",
		})
	}

	if err := models.CreateCity(&city); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "internal server error",
			"statusCode": fiber.StatusInternalServerError,
			"message":    err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":     "created",
		"statusCode": fiber.StatusCreated,
		"message":    "City has been created successfully",
		"data":       fiber.Map{"id": city.ID},
	})
}

func UpdateCity(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Invalid ID format",
		})
	}

	var updatedCity models.City

	if err := c.BodyParser(&updatedCity); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Invalid request body",
		})
	}

	rowsAffected, err := models.UpdateCityById(id, updatedCity)
	if err != nil || rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    fmt.Sprintf("City with ID %d is not found", id),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     "success",
		"statusCode": 200,
		"message":    fmt.Sprintf("City with ID %d has been updated successfully", id),
	})
}

func DeleteCity(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Invalid ID format",
		})
	}

	rowsAffected, err := models.DeleteCityById(id)
	if err != nil || rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    fmt.Sprintf("City with ID %d is not found", id),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     "success",
		"statusCode": 200,
		"message":    fmt.Sprintf("City with ID %d has been deleted successfully", id),
	})
}
