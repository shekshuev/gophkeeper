package service

import (
	"context"

	"go.uber.org/zap"

	"github.com/shekshuev/gophkeeper/internal/config"
	"github.com/shekshuev/gophkeeper/internal/logger"
	"github.com/shekshuev/gophkeeper/internal/models"
	"github.com/shekshuev/gophkeeper/internal/repository"
)

// UserServiceImpl — реализация интерфейса UserService.
// Отвечает за бизнес-логику, связанную с пользователями.
type UserServiceImpl struct {
	repo   repository.UserRepository // Репозиторий пользователей
	cfg    *config.Config            // Конфигурация приложения
	logger *logger.Logger            // Логгер
}

// NewUserServiceImpl создаёт новый экземпляр UserServiceImpl с указанным репозиторием и конфигурацией.
func NewUserServiceImpl(repo repository.UserRepository, cfg *config.Config) *UserServiceImpl {
	return &UserServiceImpl{
		repo:   repo,
		cfg:    cfg,
		logger: logger.NewLogger(),
	}
}

// GetUserByID возвращает информацию о пользователе по его идентификатору.
// Если пользователь не найден, возвращается ошибка.
func (s *UserServiceImpl) GetUserByID(ctx context.Context, id uint64) (*models.ReadUserDTO, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		s.logger.Log.Error("Ошибка при получении пользователя по ID", zap.Uint64("user_id", id), zap.Error(err))
		return nil, err
	}
	s.logger.Log.Info("Пользователь успешно получен", zap.Uint64("user_id", user.ID))
	return user, nil
}
