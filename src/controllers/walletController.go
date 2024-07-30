package controllers

import (
	// "ankasa-be/src/middlewares"
	"ankasa-be/src/middlewares"
	"ankasa-be/src/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetAllWallet(c *fiber.Ctx) error {
	results, err := models.GetAllWallet()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": results,
	})

}

func CreateWallet(c *fiber.Ctx) error {
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

	var input models.Wallet
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Invalid request body",
		})
	}
	newWalet := models.Wallet{
		CustomerID: uint(user_id),
		Saldo:      0,
	}
	if err := models.CreateWallet(&newWalet); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "Successfully created wallet",
		"statusCode": c.Status(fiber.StatusOK),
	})

}

func UpdateWallet(c *fiber.Ctx) error {
	var newWalet models.Wallet
	var wallet models.Wallet
	id, _ := strconv.Atoi(c.Params("id"))
	res := models.GetWalletById(uint(id))
	if res.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Wallet not found",
		})
	}
	if err := c.BodyParser(&newWalet); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	if err := models.UpdateWallet(uint(id), &newWalet); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update wallet",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    wallet,
		"message": "Sucessfully updated wallet",
	})
}

func DeleteWallet(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := models.DeleteWallet(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete wallet",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Wallet deleted successfully",
	})
}
