package client

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/shekshuev/gophkeeper/internal/config"
	"github.com/shekshuev/gophkeeper/internal/utils"
)

const tokenFileName = "token.json"

// tokenPath возвращает путь к файлу токена (например, ~/.gophkeeper/token.json).
func tokenPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".gophkeeper", tokenFileName)
}

// SaveToken сохраняет JWT access-токен в файл.
func SaveToken(token string) error {
	_ = os.MkdirAll(filepath.Dir(tokenPath()), 0700)
	data := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}
	b, _ := json.Marshal(data)
	return os.WriteFile(tokenPath(), b, 0600)
}

// LoadToken загружает JWT access-токен из файла.
func LoadToken() (string, error) {
	data, err := os.ReadFile(tokenPath())
	if err != nil {
		return "", err
	}
	var parsed struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return "", err
	}
	return parsed.Token, nil
}

// GetUserIDFromToken возвращает subject из access токена.
func GetUserIDFromToken() (string, error) {
	cfg := config.GetConfig()
	tokenStr, err := LoadToken()
	if err != nil {
		return "", err
	}
	claims, err := utils.GetToken(tokenStr, cfg.AccessTokenSecret)
	if err != nil {
		return "", err
	}
	if claims.Subject == "" {
		return "", fmt.Errorf("token subject is empty")
	}
	return claims.Subject, nil
}

// Logout удаляет токен с диска.
func Logout() error {
	return os.Remove(tokenPath())
}
