package controllers

import (
	"ankasa-be/src/middlewares"
	"ankasa-be/src/models"
	"ankasa-be/src/services"
	"strconv"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func WebsocketChat(c *websocket.Conn) {
	hub := c.Locals("hub").(*services.Hub)

	chat_id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		c.WriteMessage(websocket.TextMessage, []byte(`{"status":"bad request", "message":"Invalid ID format"}`))
		_ = c.Close()
		return
	}

	chat := models.SelectChatbyID(&chat_id)
	if chat.ID == 0 {
		_ = c.WriteMessage(websocket.CloseMessage, []byte{})
		_ = c.Close()
		return
	}

	room := hub.GetRoom(int(chat_id))
	room.Register <- c
	defer func() {
		room.Unregister <- c
		_ = c.Close()
	}()

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			c.WriteMessage(websocket.TextMessage, []byte(`{"status":"bad request", "message":"Failed to read message"}`))
			break
		}

		user_id, err := middlewares.JWTWebsocketAuthorize(c)
		if err != nil {
			if fiberErr, ok := err.(*fiber.Error); ok {
				c.WriteMessage(websocket.TextMessage, []byte(`{"status":"`+fiberErr.Message+`", "message":"`+fiberErr.Message+`"}`))
			} else {
				c.WriteMessage(websocket.TextMessage, []byte(`{"status":"Internal Server Error", "message":"`+err.Error()+`"}`))
			}
			break
		}

		chat := models.SelectChatUserbyChatUserID(chat_id, int(user_id))
		if chat.ID == 0 {
			c.WriteMessage(websocket.TextMessage, []byte(`{"status":"bad request", "message":"You're not allowed in this room"}`))
			break
		}

		message := models.Message{
			ChatID: uint(chat_id),
			UserID: uint(user_id),
			Body:   string(msg),
			Status: "NONE",
		}

		if err := models.CreateMessage(&message); err != nil {
			c.WriteMessage(websocket.TextMessage, []byte(`{"status":"bad request", "message":"Failed to save message"}`))
			break
		}

		room.Broadcast <- message
	}
}

func GetChatRooms(c *fiber.Ctx) error {
	user_id, err := middlewares.JWTAuthorize(c, "")
	if err != nil {
		if fiberErr, ok := err.(*fiber.Error); ok {
			return c.Status(fiberErr.Code).JSON(fiber.Map{
				"status":     fiberErr.Message,
				"statusCode": fiberErr.Code,
				"message":    fiberErr.Message,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "Internal Server Error",
			"statusCode": fiber.StatusInternalServerError,
			"message":    err.Error(),
		})
	}

	chats := models.SelectChatsbyUserID(int(user_id))
	if len(chats) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":     "no content",
			"statusCode": 202,
			"message":    "No chats available",
		})
	}

	resultChats := make([]map[string]interface{}, len(chats))
	for i, chat := range chats {

		resultMembers := make([]map[string]interface{}, len(chat.ChatUsers))
		for j, member := range chat.ChatUsers {
			customer := models.SelectCustomerfromUserID(int(member.User.ID))
			resultMembers[j] = map[string]interface{}{
				"id":           member.ID,
				"created_at":   member.CreatedAt,
				"updated_at":   member.UpdatedAt,
				"user_id":      member.User.ID,
				"customer_id":  customer.ID,
				"username":     customer.Username,
				"phone_number": customer.PhoneNumber,
				"image":        customer.Image,
			}
		}

		customer := models.SelectCustomerfromUserID(int(chat.Messages[0].User.ID))
		resultLastMessage := map[string]interface{}{
			"id":           chat.Messages[0].ID,
			"created_at":   chat.Messages[0].CreatedAt,
			"updated_at":   chat.Messages[0].UpdatedAt,
			"deleted_at":   chat.Messages[0].DeletedAt,
			"user_id":      chat.Messages[0].User.ID,
			"customer_id":  customer.ID,
			"username":     customer.Username,
			"phone_number": customer.PhoneNumber,
			"image":        customer.Image,
			"body":         chat.Messages[0].Body,
			"status":       chat.Messages[0].Status,
			"seen_at":      chat.Messages[0].SeenAt,
		}

		resultChats[i] = map[string]interface{}{
			"id":           chat.ID,
			"created_at":   chat.CreatedAt,
			"updated_at":   chat.UpdatedAt,
			"members":      resultMembers,
			"last_message": resultLastMessage,
			"status":       chat.Status,
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     "success",
		"statusCode": 200,
		"message":    "Chats OK",
		"data":       resultChats,
	})
}

func GetChatRoomById(c *fiber.Ctx) error {
	user_id, err := middlewares.JWTAuthorize(c, "")
	if err != nil {
		if fiberErr, ok := err.(*fiber.Error); ok {
			return c.Status(fiberErr.Code).JSON(fiber.Map{
				"status":     fiberErr.Message,
				"statusCode": fiberErr.Code,
				"message":    fiberErr.Message,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "Internal Server Error",
			"statusCode": fiber.StatusInternalServerError,
			"message":    err.Error(),
		})
	}

	chat_id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Invalid ID format",
		})
	}

	chat := models.SelectChatbyIDs(chat_id, int(user_id))
	if chat.ID == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":     "no content",
			"statusCode": 202,
			"message":    "Chat is empty",
		})
	}

	resultMembers := make([]map[string]interface{}, len(chat.ChatUsers))
	for i, member := range chat.ChatUsers {
		customer := models.SelectCustomerfromUserID(int(member.User.ID))
		resultMembers[i] = map[string]interface{}{
			"id":           member.ID,
			"created_at":   member.CreatedAt,
			"updated_at":   member.UpdatedAt,
			"user_id":      member.User.ID,
			"customer_id":  customer.ID,
			"username":     customer.Username,
			"phone_number": customer.PhoneNumber,
			"image":        customer.Image,
		}
	}

	resultMessages := make([]map[string]interface{}, len(chat.Messages))
	for j, message := range chat.Messages {
		customer := models.SelectCustomerfromUserID(int(message.User.ID))
		resultMessages[j] = map[string]interface{}{
			"id":           message.ID,
			"created_at":   message.CreatedAt,
			"updated_at":   message.UpdatedAt,
			"deleted_at":   message.DeletedAt,
			"user_id":      message.User.ID,
			"customer_id":  customer.ID,
			"username":     customer.Username,
			"phone_number": customer.PhoneNumber,
			"image":        customer.Image,
			"body":         message.Body,
			"status":       message.Status,
			"seen_at":      message.SeenAt,
		}
	}

	resultChats := map[string]interface{}{
		"id":         chat.ID,
		"created_at": chat.CreatedAt,
		"updated_at": chat.UpdatedAt,
		"members":    resultMembers,
		"messages":   resultMessages,
		"status":     chat.Status,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     "success",
		"statusCode": 200,
		"message":    "Chats OK",
		"data":       resultChats,
	})
}

func CreateRoom(c *fiber.Ctx) error {
	_, err := middlewares.JWTAuthorize(c, "")
	if err != nil {
		if fiberErr, ok := err.(*fiber.Error); ok {
			return c.Status(fiberErr.Code).JSON(fiber.Map{
				"status":     fiberErr.Message,
				"statusCode": fiberErr.Code,
				"message":    fiberErr.Message,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "Internal Server Error",
			"statusCode": fiber.StatusInternalServerError,
			"message":    err.Error(),
		})
	}

	var bodyRequest models.Chat
	if err := c.BodyParser(&bodyRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Invalid request body",
		})
	}

	chat := middlewares.XSSMiddleware(&bodyRequest).(*models.Chat)
	// if authErrors := helpers.StructValidation(chat); len(authErrors) > 0 {
	// 	return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
	// 		"status":     "unprocessable entity",
	// 		"statusCode": 422,
	// 		"message":    "Validation failed",
	// 		"errors":     authErrors,
	// 	})
	// }

	if err := models.CreateChat(chat); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Failed to create chat room",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":     "success",
		"statusCode": 200,
		"message":    "Chat room created successfully.",
	})
}
