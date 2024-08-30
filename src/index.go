package src

import (
	"ankasa-be/src/configs"
	"ankasa-be/src/helpers"
	"ankasa-be/src/routes"
	"ankasa-be/src/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
)

func App() *fiber.App {
	app := fiber.New()

	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:  "*",
		AllowMethods:  "GET,POST,PUT,DELETE",
		AllowHeaders:  "*",
		ExposeHeaders: "Content-Length",
	}))

	services.InitHub()
	configs.InitDB()
	helpers.Migration()
	routes.SetupRoutes(app)

	return app
}
