package service

import (
	"context"

	"github.com/shekshuev/gophkeeper/internal/config"
	"github.com/shekshuev/gophkeeper/internal/models"
	"github.com/shekshuev/gophkeeper/internal/repository"
)

// UserServiceImpl — реализация интерфейса UserService.
// Отвечает за бизнес-логику, связанную с пользователями.
type UserServiceImpl struct {
	repo repository.UserRepository // Репозиторий пользователей
	cfg  *config.Config            // Конфигурация приложения
}

// NewUserServiceImpl создаёт новый экземпляр UserServiceImpl с указанным репозиторием и конфигурацией.
func NewUserServiceImpl(repo repository.UserRepository, cfg *config.Config) *UserServiceImpl {
	return &UserServiceImpl{repo: repo, cfg: cfg}
}

// GetUserByID возвращает информацию о пользователе по его идентификатору.
// Если пользователь не найден, возвращается ошибка.
func (s *UserServiceImpl) GetUserByID(ctx context.Context, id uint64) (*models.ReadUserDTO, error) {
	return s.repo.GetUserByID(ctx, id)
}
