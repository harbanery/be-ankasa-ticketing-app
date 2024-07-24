package controllers

import (
	// "ankasa-be/src/middlewares"
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

// func CreateWallet(c *fiber.Ctx) error {
// 	var input models.Wallet
// 	if err := c.BodyParser(&input); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"status":     "bad request",
// 			"statusCode": 400,
// 			"message":    "Invalid request body",
// 		})
// 	}
// 	wallet := middlewares.XSSMiddleware(&input).(*models.Wallet)
// newWalet := models.Wallet{
// 	CustomerID: ,
// 	}
// }

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
