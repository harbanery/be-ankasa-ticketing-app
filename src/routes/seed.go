package routes

import (
	"ankasa-be/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func seedRoutes(app *fiber.App) {
	seeds := app.Group("/seeds")

	seeds.Get("/merchants", controllers.GenerateMerchantSeed)
	seeds.Get("/destinations", controllers.GenerateCityCountrySeed)
}
