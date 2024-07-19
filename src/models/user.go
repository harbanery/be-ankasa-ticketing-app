package models

import (
	"ankasa-be/src/configs"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role" validate:"oneof=merchant customer"`
	IsVerify string `json:"is_verify"`
}

type Register struct {
	gorm.Model
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Username string `json:"username" validate:"required,max=50"`
	Role     string `json:"role" validate:"oneof=merchant customer"`
}

// type Merchant struct {
// 	gorm.Model
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// 	Username string `json:"username"`
// 	Image    string `json:"image"`
// }

type UserVerification struct {
	gorm.Model
	UserID uint   `json:"user_id"`
	Token  string `json:"token"`
}

func SelectUsers() []*User {
	var users []*User
	configs.DB.Find(&users)
	return users
}

func SelectUserfromID(id int) *User {
	var user *User
	configs.DB.First(&user, "id = ?", id)
	return user
}

func SelectUserfromEmail(email string) *User {
	var user *User
	configs.DB.First(&user, "email = ?", email)
	return user
}

func SelectUserVerification(user_id int, token string) *UserVerification {
	var userVerification UserVerification
	configs.DB.Where("user_id = ? AND token = ?", user_id, token).First(&userVerification)
	return &userVerification
}

func CreateUser(user *User) (uint, error) {
	result := configs.DB.Create(&user)
	return user.ID, result.Error
}

func CreateUserVerification(userVerification *UserVerification) error {
	result := configs.DB.Create(&userVerification)
	return result.Error
}

func UpdateUserfromEmail(email string, updatedUser *User) error {
	result := configs.DB.Model(&User{}).Where("email = ?", email).Updates(updatedUser)
	return result.Error
}

func UpdateUserVerify(id int) error {
	result := configs.DB.Model(&User{}).Where("id = ?", id).Update("is_verify", "true")
	return result.Error
}

func DeleteUserVerification(id int, token string) error {
	result := configs.DB.Where("id = ? AND token = ?", id, token).Delete(&UserVerification{})
	return result.Error
}
