package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword принимает обычный пароль и возвращает его bcrypt-хеш.
// В случае ошибки возвращает пустую строку.
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(bytes)
}

// VerifyPassword сравнивает исходный пароль с ранее сохранённым хешем.
// Возвращает true, если пароль корректен, иначе false.
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
