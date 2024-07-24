package routes

import (
	"ankasa-be/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func walletRoutes(app *fiber.App) {
	wallet := app.Group("/wallets")
	wallet.Get("/", controllers.GetAllWallet)
	wallet.Put("/update/:id", controllers.UpdateWallet)
	wallet.Delete("/delete/:id", controllers.DeleteWallet)
}
