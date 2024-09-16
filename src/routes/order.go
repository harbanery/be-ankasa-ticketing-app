package routes

import (
	"ankasa-be/src/controllers"
	"ankasa-be/src/middlewares"

	"github.com/gofiber/fiber/v2"
)

func ordersRoutes(app *fiber.App) {
	orders := app.Group("/orders")
	orders.Get("/", controllers.GetAllOrders)
	// orders.Get("/:id", controllers.GetTicketById)
	orders.Post("/", middlewares.JWTMiddleware(), controllers.CreateOrder)
	// orders.Put("/:id", controllers.UpdateTicket)
	// orders.Delete("/:id", controllers.DeleteTicket)
}
