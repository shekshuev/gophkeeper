package service

import (
	"context"

	"go.uber.org/zap"

	"github.com/shekshuev/gophkeeper/internal/logger"
	"github.com/shekshuev/gophkeeper/internal/models"
	"github.com/shekshuev/gophkeeper/internal/repository"
)

// SecretServiceImpl реализует SecretService.
// Отвечает за бизнес-логику по работе с пользовательскими секретами.
type SecretServiceImpl struct {
	repo   repository.SecretRepository // Репозиторий секретов
	logger *logger.Logger              // Логгер
}

// NewSecretServiceImpl создаёт новый экземпляр сервиса секретов.
func NewSecretServiceImpl(repo repository.SecretRepository) *SecretServiceImpl {
	return &SecretServiceImpl{
		repo:   repo,
		logger: logger.NewLogger(),
	}
}

// Create сохраняет новый секрет.
// Возвращает ID созданного секрета или ошибку.
func (s *SecretServiceImpl) Create(ctx context.Context, dto models.CreateSecretDTO) (uint64, error) {
	id, err := s.repo.Create(ctx, dto)
	if err != nil {
		s.logger.Log.Error("Не удалось сохранить секрет", zap.Uint64("user_id", dto.UserID), zap.String("title", dto.Title), zap.Error(err))
		return 0, err
	}
	s.logger.Log.Info("Секрет успешно создан", zap.Uint64("secret_id", id), zap.Uint64("user_id", dto.UserID))
	return id, nil
}

// GetByID возвращает секрет по ID.
// Если секрет не найден, возвращает nil, nil.
func (s *SecretServiceImpl) GetByID(ctx context.Context, id uint64) (*models.ReadSecretDTO, error) {
	secret, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Log.Error("Ошибка при получении секрета по ID", zap.Uint64("secret_id", id), zap.Error(err))
		return nil, err
	}
	if secret == nil {
		s.logger.Log.Warn("Секрет не найден", zap.Uint64("secret_id", id))
	} else {
		s.logger.Log.Info("Секрет успешно получен", zap.Uint64("secret_id", secret.ID))
	}
	return secret, nil
}

// GetAllByUser возвращает все секреты конкретного пользователя.
func (s *SecretServiceImpl) GetAllByUser(ctx context.Context, userID uint64) ([]models.ReadSecretDTO, error) {
	secrets, err := s.repo.GetAllByUser(ctx, userID)
	if err != nil {
		s.logger.Log.Error("Ошибка при получении секретов пользователя", zap.Uint64("user_id", userID), zap.Error(err))
		return nil, err
	}
	s.logger.Log.Info("Секреты пользователя успешно получены", zap.Uint64("user_id", userID), zap.Int("count", len(secrets)))
	return secrets, nil
}

// DeleteByID удаляет секрет по ID.
// Возвращает ошибку, если удаление не удалось.
func (s *SecretServiceImpl) DeleteByID(ctx context.Context, id uint64) error {
	err := s.repo.DeleteByID(ctx, id)
	if err != nil {
		s.logger.Log.Error("Ошибка при удалении секрета", zap.Uint64("secret_id", id), zap.Error(err))
		return err
	}
	s.logger.Log.Info("Секрет успешно удалён", zap.Uint64("secret_id", id))
	return nil
}
