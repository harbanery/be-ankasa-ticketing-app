package routes

import (
	"ankasa-be/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func countriesRoutes(app *fiber.App) {
	countries := app.Group("/countries")
	countries.Get("/", controllers.GetAllCountries)
	countries.Get("/:id", controllers.GetCountryById)
	countries.Post("/", controllers.CreateCountry)
	countries.Put("/:id", controllers.UpdateCountry)
	countries.Delete("/:id", controllers.DeleteCountry)
}
