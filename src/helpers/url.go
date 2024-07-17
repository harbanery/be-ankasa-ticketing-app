package helpers

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"strconv"
)

func GenerateURL(userID int, pathname string) (string, string, error) {
	uuid := strconv.Itoa(userID)

	tokenBytes := make([]byte, 64)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", "", err
	}

	token := hex.EncodeToString(tokenBytes)

	url := os.Getenv("API_URL") + pathname + "?id=" + uuid + "&token=" + token
	return url, token, nil
}
