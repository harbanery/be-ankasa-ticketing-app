package routes

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	mainRoutes(app)
	userRoutes(app)
	merchantRoutes(app)
	customerRoutes(app)
	ticketsRoutes(app)
	categoriesRoutes(app)
	countriesRoutes(app)
	citiesRoutes(app)
	chatsRoutes(app)
	seedRoutes(app)
	walletRoutes(app)
}
