package routes

import (
	"ankasa-be/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func categoriesRoutes(app *fiber.App) {
	categories := app.Group("/categories")
	categories.Get("/", controllers.GetAllCategories)
	categories.Get("/:id", controllers.GetCategoryById)
	categories.Post("/", controllers.CreateCategory)
	categories.Put("/:id", controllers.UpdateCategory)
	categories.Delete("/:id", controllers.DeleteCategory)
}
