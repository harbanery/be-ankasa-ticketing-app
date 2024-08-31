package routes

import (
	"ankasa-be/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func merchantRoutes(app *fiber.App) {
	customer := app.Group("/merchant")
	customer.Get("/seed", controllers.GenerateMerchantSeed)
}
