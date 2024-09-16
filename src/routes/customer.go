package routes

import (
	"ankasa-be/src/controllers"
	"ankasa-be/src/middlewares"

	"github.com/gofiber/fiber/v2"
)

func customerRoutes(app *fiber.App) {
	customer := app.Group("/customer")
	// Buat testing aja
	customer.Get("/", controllers.GetCustomers)
	// 
	
	customer.Get("/profile", middlewares.JWTMiddleware(), controllers.GetCustomerProfile)
	customer.Put("/profile", middlewares.JWTMiddleware(), controllers.UpdateCustomerProfile)
	customer.Post("/upload", middlewares.JWTMiddleware(), controllers.UploadProfileImage)
}
