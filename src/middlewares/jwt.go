package middlewares

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func ExtractToken(c *fiber.Ctx) string {
	bearerToken := c.Get("Authorization")
	if strings.HasPrefix(bearerToken, "Bearer ") {
		return strings.TrimPrefix(bearerToken, "Bearer ")
	}
	return ""
}

func JWTMiddleware() fiber.Handler {
	secretKey := os.Getenv("SECRETKEY")
	return func(c *fiber.Ctx) error {
		tokenString := ExtractToken(c)
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":     "unauthorized",
				"statusCode": 401,
				"message":    "Token string unauthorized",
			})
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":     "unauthorized",
				"statusCode": 401,
				"message":    "Token unauthorized",
			})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Locals("user", claims)
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":     "unauthorized",
				"statusCode": 401,
				"message":    "Token claims unauthorized",
			})
		}

		return c.Next()
	}
}

func JWTAuthorize(c *fiber.Ctx, requiredRole string) (float64, error) {
	user, ok := c.Locals("user").(jwt.MapClaims)
	if !ok {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	if requiredRole != "" {
		role, ok := user["role"].(string)
		if !ok || role != requiredRole {
			return 0, fiber.NewError(fiber.StatusForbidden, "Incorrect role")
		}
	}

	id, ok := user["id"].(float64)
	if !ok {
		return 0, fiber.NewError(fiber.StatusBadRequest, "Invalid ID format")
	}

	return id, nil
}

func JWTWebsocketAuthorize(c *websocket.Conn) (float64, error) {
	user, ok := c.Locals("user").(jwt.MapClaims)
	if !ok {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	id, ok := user["id"].(float64)
	if !ok {
		return 0, fiber.NewError(fiber.StatusBadRequest, "Invalid ID format")
	}

	return id, nil
}
