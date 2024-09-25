package routes

import (
	"ankasa-be/src/controllers"
	"ankasa-be/src/middlewares"

	"github.com/gofiber/fiber/v2"
)

func ordersRoutes(app *fiber.App) {
	orders := app.Group("/orders")
	orders.Get("/", middlewares.JWTMiddleware(), controllers.GetAllOrders)
	orders.Get("/:id", middlewares.JWTMiddleware(), controllers.GetBookingPass)
	orders.Post("/", middlewares.JWTMiddleware(), controllers.CreatePaymentOrder)
	orders.Post("/payment_method", controllers.HandlePaymentMethodCallback)
	orders.Post("/payment", controllers.HandlePaymentRequestCallback)
}
