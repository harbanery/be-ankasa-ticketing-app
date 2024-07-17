package controllers

import (
	"ankasa-be/src/helpers"
	"ankasa-be/src/middlewares"
	"ankasa-be/src/models"
	"ankasa-be/src/services"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func GetUsers(c *fiber.Ctx) error {
	users := models.SelectUsers()
	if len(users) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    "Users not found.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     "success",
		"statusCode": 200,
		"data":       users,
	})
}

func RegisterUser(c *fiber.Ctx) error {
	var register models.Register
	if err := c.BodyParser(&register); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Invalid request body",
		})
	}

	user := middlewares.XSSMiddleware(&register).(*models.Register)
	if authErrors := helpers.PasswordValidation(user.Password, helpers.StructValidation(user)); len(authErrors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":     "unprocessable entity",
			"statusCode": 422,
			"message":    "Validation failed",
			"errors":     authErrors,
		})
	}

	if existUser := models.SelectUserfromEmail(user.Email); existUser.ID != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Email already exists",
		})
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Password error",
		})
	}

	newUser := models.User{
		Email:    user.Email,
		Password: string(hashPassword),
		Role:     user.Role,
		IsVerify: "false",
	}

	userID, err := models.CreateUser(&newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Failed to create user",
		})
	}

	url, token, err := helpers.GenerateURL(int(userID), "verify")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Failed to generate URL",
		})
	}

	if err := services.SendEmail(newUser.Email, "Verify Email", url); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Failed to send verification email",
		})
	}

	newUserVerification := models.UserVerification{
		UserID: userID,
		Token:  token,
	}

	if err := models.CreateUserVerification(&newUserVerification); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Failed to create user verification",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":     "success",
		"statusCode": 200,
		"message":    "User created successfully. Please check in your email to verify.",
	})
}

func VerificationAccount(c *fiber.Ctx) error {
	queryUserID := c.Query("id")
	queryToken := c.Query("token")

	if queryUserID == "" || queryToken == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":     "forbidden",
			"statusCode": 403,
			"error":      "Invalid url verification",
		})
	}

	userID, err := strconv.Atoi(queryUserID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":     "forbidden",
			"statusCode": 403,
			"error":      "Invalid user ID",
		})
	}

	existUser := models.SelectUserfromID(userID)
	if existUser.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    "Email not found",
		})
	}

	if existUser.IsVerify != "false" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":     "forbidden",
			"statusCode": 403,
			"message":    "Users has been verified",
		})
	}

	existUserVerify := models.SelectUserVerification(userID, queryToken)
	if existUserVerify.ID == 0 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":     "forbidden",
			"statusCode": 403,
			"message":    "Error invalid credential verification",
		})
	}

	if err := models.UpdateUserVerify(userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Failed to update account verification",
		})
	}

	if err := models.DeleteUserVerification(int(existUserVerify.ID), queryToken); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Failed to delete user verification",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     "ok",
		"statusCode": 200,
		"message":    "Users verified successfully",
	})
}

func CreateRefreshToken(c *fiber.Ctx) error {
	var refreshData struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&refreshData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Failed to parse request body",
		})
	}

	token, err := helpers.GenerateToken(os.Getenv("SECRETKEY"), map[string]interface{}{"refreshToken": refreshData.RefreshToken})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Could not generate access token",
		})
	}

	refreshToken, err := helpers.GenerateRefreshToken(os.Getenv("SECRETKEY"), map[string]interface{}{"refreshToken": refreshData.RefreshToken})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Could not generate refresh token",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":        "success",
		"statusCode":    201,
		"message":       "Refresh successfully",
		"token":         token,
		"refresh_token": refreshToken,
	})
}

func LogoutUser(c *fiber.Ctx) error {
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":     "success",
		"statusCode": 201,
		"message":    "Logout successfully",
	})
}

func RequestResetPassword(c *fiber.Ctx) error {
	var requestEmail models.User
	if err := c.BodyParser(&requestEmail); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Invalid request body",
		})
	}

	user := middlewares.XSSMiddleware(&requestEmail).(*models.User)
	if authErrors := helpers.FieldRequiredValidation(user.Email, "required,email"); authErrors != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":     "unprocessable entity",
			"statusCode": 422,
			"message":    "Validation failed",
			"errors":     authErrors,
		})
	}

	existUser := models.SelectUserfromEmail(user.Email)
	if existUser.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    "Email not found",
		})
	}

	url, _, err := helpers.GenerateURL(int(existUser.ID), "resetPassword")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Failed to generate URL",
		})
	}

	if err := services.SendEmail(user.Email, "Reset Password", url); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Failed to send reset Password",
		})
	}

	// newUserVerification := models.UserVerification{
	// 	UserID: existUser.ID,
	// 	Token:  token,
	// }

	// if err := models.CreateUserVerification(&newUserVerification); err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"status":     "server error",
	// 		"statusCode": 500,
	// 		"message":    "Failed to create user verification",
	// 	})
	// }

	if err := models.UpdateUserVerificationfromID(int(existUser.ID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Failed to create user verification",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"status":      "accepted",
		"statusCode":  202,
		"message":     "Password reset email sent. Please check in your email to reset your password.",
		"expectedUrl": url,
	})
}

func ResetPassword(c *fiber.Ctx) error {
	queryUserId := c.Query("id")
	queryToken := c.Query("token")

	if queryUserId == "" || queryToken == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Invalid url verification",
		})
	}

	userID, err := strconv.Atoi(queryUserId)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	existUser := models.SelectUserfromID(userID)
	if existUser.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    "Email not found",
		})
	}

	var updatedUser models.User
	if err := c.BodyParser(&updatedUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Invalid request body",
		})
	}
	updatedUser.Email = existUser.Email

	user := middlewares.XSSMiddleware(&updatedUser).(*models.User)
	if authErrors := helpers.PasswordValidation(user.Password, helpers.StructValidation(user)); len(authErrors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":     "unprocessable entity",
			"statusCode": 422,
			"message":    "Validation failed",
			"errors":     authErrors,
		})
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Failed to hash password",
		})
	}
	user.Password = string(hashPassword)

	existUserVerify := models.SelectUserVerification(userID, queryToken)
	if existUserVerify.ID == 0 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":     "forbidden",
			"statusCode": 403,
			"message":    "Error invalid credential verification",
		})
	}

	if err := models.UpdateUserfromEmail(user.Email, user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Failed to reset password",
		})
	}

	if err := models.DeleteUserVerification(int(existUserVerify.ID), queryToken); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Failed to delete user verification",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":     "success",
		"statusCode": 201,
		"message":    "Reset password successfully",
	})
}
