package routes

import (
	"ankasa-be/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func ticketsRoutes(app *fiber.App) {
	tickets := app.Group("/tickets")
	tickets.Get("/", controllers.GetAllTickets)
	tickets.Get("/:id", controllers.GetTicketById)
	tickets.Post("/", controllers.CreateTicket)
	tickets.Put("/:id", controllers.UpdateTicket)
	tickets.Delete("/:id", controllers.DeleteTicket)
}
