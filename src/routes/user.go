package routes

import (
	"ankasa-be/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func userRoutes(app *fiber.App) {
	app.Get("/user", controllers.GetUsers)
	app.Post("/register", controllers.RegisterUser)
	app.Get("/verify", controllers.VerificationAccount)
	app.Post("/loginAuthProvider", controllers.LoginUserwithAuthProvider)
	app.Post("/login", controllers.LoginUser)
	app.Get("/logout", controllers.LogoutUser)
	app.Post("/requestResetPassword", controllers.RequestResetPassword)
	app.Put("/resetPassword", controllers.ResetPassword)
	app.Post("/refreshToken", controllers.CreateRefreshToken)
}
