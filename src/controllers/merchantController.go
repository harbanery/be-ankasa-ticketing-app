package controllers

import (
	"ankasa-be/src/models"

	"github.com/gofiber/fiber/v2"
)

func GetAllMerchants(c *fiber.Ctx) error {
	merchants := models.SelectAllMerchants()
	if len(merchants) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    "Merchants not found",
		})
	}

	resultMerchants := make([]map[string]interface{}, len(merchants))
	for i, merchant := range merchants {

		resultClasses := make([]map[string]interface{}, len(merchant.Classes))
		for j, class := range merchant.Classes {
			resultClasses[j] = map[string]interface{}{
				"id":               class.ID,
				"created_at":       class.CreatedAt,
				"updated_at":       class.UpdatedAt,
				"name":             class.Name,
				"price":            class.Price,
				"seats":            class.Seats,
				"row_seats":        class.RowSeats,
				"is_refund":        class.IsRefund,
				"is_reschedule":    class.IsReschedule,
				"is_luggage":       class.IsLuggage,
				"is_inflight_meal": class.IsInflightMeal,
				"is_wifi":          class.IsWifi,
			}
		}

		resultMerchants[i] = map[string]interface{}{
			"id":          merchant.ID,
			"created_at":  merchant.CreatedAt,
			"updated_at":  merchant.UpdatedAt,
			"name":        merchant.Name,
			"image":       merchant.Image,
			"description": merchant.Description,
			"classes":     resultClasses,
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     "success",
		"statusCode": 200,
		"message":    "Chats OK",
		"data":       resultMerchants,
	})
}
