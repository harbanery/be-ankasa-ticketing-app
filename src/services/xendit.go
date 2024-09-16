package services

import (
	"os"

	"github.com/xendit/xendit-go/v6"
)

var Client *xendit.APIClient

func InitXendit() {
	secretKey := os.Getenv("XENDIT_SECRET_KEY")

	Client = xendit.NewClient(secretKey)
}
