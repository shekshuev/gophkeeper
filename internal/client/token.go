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

// tokenPath — возвращает абсолютный путь к файлу токена.
// Например: ~/.gophkeeper/token.json.
func tokenPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".gophkeeper", tokenFileName)
}

// SaveToken — сохраняет JWT access-токен в файл.
// Создаёт директорию ~/.gophkeeper при необходимости.
// Токен сохраняется в JSON-формате с правами 0600.
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

// LoadToken — загружает access-токен из локального файла.
// Возвращает строку токена или ошибку, если файл не существует или повреждён.
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

// GetUserIDFromToken — парсит токен и возвращает значение поля subject (user ID).
// Используется для определения текущего пользователя в клиенте.
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

// Logout — удаляет локальный файл токена (выход из системы).
func Logout() error {
	return os.Remove(tokenPath())
}
