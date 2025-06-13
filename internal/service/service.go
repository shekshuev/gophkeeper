package service

import (
	"context"
	"fmt"

	"github.com/shekshuev/gophkeeper/internal/models"
)

// UserService определяет операции для управления пользователями и получения информации о них.
type UserService interface {
	// GetAllUsers возвращает список пользователей с пагинацией.
	// limit — максимальное количество пользователей,
	// offset — смещение от начала списка.
	GetAllUsers(ctx context.Context, limit, offset uint64) ([]models.ReadUserDTO, error)

	// GetUserByID возвращает пользователя по его уникальному идентификатору.
	GetUserByID(ctx context.Context, id uint64) (*models.ReadUserDTO, error)
}

// AuthService отвечает за регистрацию и аутентификацию пользователей.
type AuthService interface {
	// Login проверяет логин и пароль пользователя.
	// При успешной аутентификации возвращает access и refresh токены.
	Login(ctx context.Context, dto models.LoginUserDTO) (*models.ReadTokenDTO, error)

	// Register создаёт нового пользователя и сразу авторизует его.
	// Возвращает access и refresh токены.
	Register(ctx context.Context, dto models.RegisterUserDTO) (*models.ReadTokenDTO, error)
}

// ErrUserNotFound возвращается, если пользователь не найден в базе.
var ErrUserNotFound = fmt.Errorf("user not found")

// ErrWrongPassword возвращается, если пароль не совпадает с сохранённым хешем.
var ErrWrongPassword = fmt.Errorf("wrong password")
