package routes

import (
	"ankasa-be/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func citiesRoutes(app *fiber.App) {
	cities := app.Group("/cities")
	cities.Get("/", controllers.GetAllCities)
	cities.Get("/:id", controllers.GetCityById)
	cities.Post("/", controllers.CreateCity)
	cities.Put("/:id", controllers.UpdateCity)
	cities.Delete("/:id", controllers.DeleteCity)
}
