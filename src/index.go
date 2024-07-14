package src

import (
	"ankasa-be/src/configs"
	"ankasa-be/src/helpers"
	"ankasa-be/src/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func App() *fiber.App {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:  "*",
		AllowMethods:  "GET,POST,PUT,DELETE",
		AllowHeaders:  "*",
		ExposeHeaders: "Content-Length",
	}))

	configs.InitDB()
	helpers.Migration()
	routes.SetupRoutes(app)

	return app
}
