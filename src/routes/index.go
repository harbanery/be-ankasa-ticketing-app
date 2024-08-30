package routes

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	mainRoutes(app)
	userRoutes(app)
	customerRoutes(app)
	ticketsRoutes(app)
	categoriesRoutes(app)
	countriesRoutes(app)
	citiesRoutes(app)
	chatsRoutes(app)
}
