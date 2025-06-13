package repository

import (
	"context"
	"fmt"

	"github.com/shekshuev/gophkeeper/internal/models"
)

// DatabaseChecker определяет поведение для проверки доступности базы данных.
type DatabaseChecker interface {
	// CheckDBConnection проверяет соединение с базой данных.
	// Возвращает ошибку в случае недоступности.
	CheckDBConnection() error
}

// UserRepository определяет поведение для репозитория пользователей.
type UserRepository interface {
	// GetUserByUserName находит пользователя по его userName.
	// Возвращает ReadAuthUserDataDTO или ошибку, если пользователь не найден.
	GetUserByUserName(ctx context.Context, userName string) (*models.ReadAuthUserDataDTO, error)

	// GetUserByUserName находит пользователя по его уникальному идентификатору.
	// Возвращает ReadUserDTO или ошибку, если пользователь не найден.
	GetUserByID(ctx context.Context, id uint64) (*models.ReadUserDTO, error)

	// CreateUser создает нового пользователя на основе CreateUserDTO.
	// Возвращает ReadAuthUserDataDTO или ошибку, если пользователь уже существует или произошла ошибка при сохранении.
	CreateUser(ctx context.Context, user models.CreateUserDTO) (*models.ReadAuthUserDataDTO, error)
}

// ErrNotFound используется, когда запись не найдена в базе данных.
var ErrNotFound = fmt.Errorf("not found")

// ErrUserExists используется, когда попытка создать пользователя с уже существующим userName.
var ErrUserExists = fmt.Errorf("user already exists")
