package repository

import (
	"context"
	"fmt"

	"github.com/shekshuev/gophkeeper/internal/models"
)

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

// SecretRepository определяет интерфейс для работы с секретами.
type SecretRepository interface {
	// Create сохраняет новый секрет для пользователя.
	// Возвращает ID созданного секрета или ошибку.
	Create(ctx context.Context, dto models.CreateSecretDTO) (uint64, error)

	// GetByID возвращает секрет по его ID.
	// Если секрет не найден, возвращается ошибка ErrNotFound.
	GetByID(ctx context.Context, id uint64) (*models.ReadSecretDTO, error)

	// GetAllByUser возвращает все секреты пользователя по его userID.
	GetAllByUser(ctx context.Context, userID uint64) ([]models.ReadSecretDTO, error)

	// DeleteByID удаляет секрет по его ID.
	// Если секрета нет — ничего не делает.
	DeleteByID(ctx context.Context, id uint64) error
}

// ErrNotFound используется, когда запись не найдена в базе данных.
var ErrNotFound = fmt.Errorf("not found")

// ErrUserExists используется, когда попытка создать пользователя с уже существующим userName.
var ErrUserExists = fmt.Errorf("user already exists")

// ErrMarshalPayload возникает при ошибке сериализации (marshal) данных секрета в JSON перед сохранением в БД.
var ErrMarshalPayload = fmt.Errorf("error marshal payload")

// ErrUnmarshalPayload возникает при ошибке десериализации (unmarshal) JSON-данных секрета, полученных из БД.
var ErrUnmarshalPayload = fmt.Errorf("error unmarshal payload")
