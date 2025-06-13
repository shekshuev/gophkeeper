package models

import "time"

// CreateUserDTO используется для передачи данных при создании нового пользователя.
// Поле PasswordHash содержит уже захешированный пароль (например, bcrypt).
type CreateUserDTO struct {
	UserName     string // Уникальное имя пользователя (логин)
	PasswordHash string // Хеш пароля
	FirstName    string // Имя пользователя
	LastName     string // Фамилия пользователя
}

// ReadUserDTO представляет данные пользователя, возвращаемые при запросах (например, в списке или профиле).
type ReadUserDTO struct {
	ID        uint64    `json:"id"`         // Уникальный идентификатор пользователя
	UserName  string    `json:"user_name"`  // Логин пользователя
	FirstName string    `json:"first_name"` // Имя
	LastName  string    `json:"last_name"`  // Фамилия
	CreatedAt time.Time `json:"created_at"` // Время создания записи
	UpdatedAt time.Time `json:"updated_at"` // Время последнего обновления
}

// ReadAuthUserDataDTO содержит минимальный набор данных для авторизации и валидации логина/пароля.
// Используется в репозиториях и сервисе аутентификации.
type ReadAuthUserDataDTO struct {
	ID           uint64 `json:"id"`            // Идентификатор пользователя
	UserName     string `json:"user_name"`     // Логин
	PasswordHash string `json:"password_hash"` // Хеш пароля (для проверки)
}
