package routes

import (
	"ankasa-be/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func userRoutes(app *fiber.App) {
	app.Get("/users", controllers.GetUsers)
	app.Post("/users", controllers.RegisterUser)
}
