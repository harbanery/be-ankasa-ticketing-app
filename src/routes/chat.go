package routes

import (
	"ankasa-be/src/controllers"
	"ankasa-be/src/middlewares"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func chatsRoutes(app *fiber.App) {
	ws := app.Group("/ws/chat")
	ws.Use("/", middlewares.AllowUpgrade)
	ws.Get("/:id", middlewares.JWTMiddleware(), websocket.New(controllers.WebsocketChat))

	chats := app.Group("/chats")
	chats.Get("/", middlewares.JWTMiddleware(), controllers.GetChatRooms)
	chats.Get("/:id", middlewares.JWTMiddleware(), controllers.GetChatRoomById)
	chats.Post("/", middlewares.JWTMiddleware(), controllers.CreateRoom)
	// categories.Put("/:id", controllers.UpdateCategory)
	// categories.Delete("/:id", controllers.DeleteCategory)
}
