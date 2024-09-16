package models

import (
	"ankasa-be/src/configs"
	"time"

	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model
	ChatUsers []ChatUser `json:"members" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Messages  []Message  `json:"messages"`
	Status    string     `json:"status"` // NONE, REPLY
}

type ChatUser struct {
	gorm.Model
	ChatID uint `json:"chat_id"`
	UserID uint `json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user"`
}

type Message struct {
	gorm.Model
	ChatID uint       `json:"chat_id"`
	UserID uint       `json:"user_id"`
	User   User       `gorm:"foreignKey:UserID" json:"user"`
	Body   string     `json:"body"`
	Status string     `gorm:"default:NONE" json:"status"` // NONE, SENT, UNSEEN, SEEN
	SeenAt *time.Time `json:"seen_at"`
}

// func CreateChat(chat *Chat) (uint, error) {
// 	result := configs.DB.Create(&chat)
// 	return chat.ID, result.Error
// }

func CreateChat(chat *Chat) error {
	result := configs.DB.Create(&chat)
	return result.Error
}

func CreateMessage(message *Message) error {
	result := configs.DB.Create(&message)
	return result.Error
}

func SelectChatbyID(id *int) *Chat {
	var chat *Chat
	configs.DB.First(&chat, "id = ?", &id)
	return chat
}

func SelectChatUserbyChatUserID(chat_id, user_id int) *ChatUser {
	var member *ChatUser
	configs.DB.First(&member, "chat_id = ? AND user_id = ?", &chat_id, &user_id)
	return member
}

func SelectChatsbyUserID(user_id int) []*Chat {
	var chats []*Chat

	configs.DB.Preload("Messages", func(db *gorm.DB) *gorm.DB {
		return db.Preload("User").Order("created_at DESC").Limit(1)
	}).Preload("ChatUsers", func(db *gorm.DB) *gorm.DB {
		return db.Preload("User").Where("user_id != ?", user_id)
	}).Joins("INNER JOIN chat_users ON chat_users.chat_id = chats.id AND chat_users.user_id IN (?)", user_id).
		Joins("INNER JOIN messages ON messages.chat_id = chats.id").
		Group("chats.id").Having("COUNT(messages.id) > 0").
		Find(&chats)

	return chats
}

func SelectChatbyIDs(chat_id, user_id int) *Chat {
	var chat *Chat

	configs.DB.Preload("Messages", func(db *gorm.DB) *gorm.DB {
		return db.Preload("User").Order("created_at DESC")
	}).Preload("ChatUsers", func(db *gorm.DB) *gorm.DB {
		return db.Preload("User").Where("user_id != ?", user_id)
	}).Joins("INNER JOIN chat_users ON chat_users.chat_id = chats.id AND chat_users.user_id IN (?)", user_id).
		Joins("INNER JOIN messages ON messages.chat_id = chats.id").
		Group("chats.id").First(&chat)

	return chat
}
