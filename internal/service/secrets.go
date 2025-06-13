package service

import (
	"context"

	"github.com/shekshuev/gophkeeper/internal/models"
	"github.com/shekshuev/gophkeeper/internal/repository"
)

// SecretServiceImpl реализует SecretService.
type SecretServiceImpl struct {
	repo repository.SecretRepository // Репозиторий секретов
}

// NewSecretService создаёт новый экземпляр сервиса секретов.
func NewSecretServiceImpl(repo repository.SecretRepository) *SecretServiceImpl {
	return &SecretServiceImpl{repo: repo}
}

// Create сохраняет новый секрет.
func (s *SecretServiceImpl) Create(ctx context.Context, dto models.CreateSecretDTO) (uint64, error) {
	return s.repo.Create(ctx, dto)
}

// GetByID возвращает секрет по ID.
func (s *SecretServiceImpl) GetByID(ctx context.Context, id uint64) (*models.ReadSecretDTO, error) {
	return s.repo.GetByID(ctx, id)
}

// GetAllByUser возвращает все секреты конкретного пользователя.
func (s *SecretServiceImpl) GetAllByUser(ctx context.Context, userID uint64) ([]models.ReadSecretDTO, error) {
	return s.repo.GetAllByUser(ctx, userID)
}

// DeleteByID удаляет секрет по ID.
func (s *SecretServiceImpl) DeleteByID(ctx context.Context, id uint64) error {
	return s.repo.DeleteByID(ctx, id)
}
