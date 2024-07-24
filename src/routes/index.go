package routes

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	mainRoutes(app)
	userRoutes(app)
	customerRoutes(app)
	walletRoutes(app)
}
