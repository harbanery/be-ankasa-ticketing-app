package models

import (
	"ankasa-be/src/configs"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SelectUsers() []*User {
	var users []*User
	configs.DB.Find(&users)
	return users
}

func CreateUser(user *User) error {
	result := configs.DB.Create(&user)
	return result.Error
}
