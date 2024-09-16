package routes

import (
	"ankasa-be/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func merchantRoutes(app *fiber.App) {
	merchants := app.Group("/merchants")
	merchants.Get("", controllers.GetAllMerchants)
}
