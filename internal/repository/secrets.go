package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
	"go.uber.org/zap"

	"github.com/shekshuev/gophkeeper/internal/config"
	"github.com/shekshuev/gophkeeper/internal/logger"
	"github.com/shekshuev/gophkeeper/internal/models"
)

// SecretRepositoryImpl — реализация интерфейса SecretRepository для работы с секретами в PostgreSQL.
type SecretRepositoryImpl struct {
	db     *sql.DB        // соединение с базой данных
	cfg    *config.Config // конфигурация приложения
	logger *logger.Logger // логгер
}

// NewSecretRepositoryImpl создаёт новый экземпляр SecretRepositoryImpl.
// Устанавливает соединение с базой данных на основе переданного DSN из конфигурации.
func NewSecretRepositoryImpl(cfg *config.Config) *SecretRepositoryImpl {
	log := logger.NewLogger()

	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		log.Log.Fatal("Не удалось подключиться к базе данных", zap.Error(err))
	}
	log.Log.Info("Установлено соединение с базой данных (secrets)")

	return &SecretRepositoryImpl{
		db:     db,
		cfg:    cfg,
		logger: log,
	}
}

// Create сохраняет новый секрет в базу данных.
// Принимает DTO с userID, названием и данными (в виде map).
// Возвращает ID созданного секрета или ошибку.
func (r *SecretRepositoryImpl) Create(ctx context.Context, dto models.CreateSecretDTO) (uint64, error) {
	dataBytes, err := json.Marshal(dto.Data)
	if err != nil {
		r.logger.Log.Error("Ошибка маршалинга данных секрета", zap.Error(err))
		return 0, ErrMarshalPayload
	}

	query := `
		insert into secrets (user_id, title, data)
		values ($1, $2, $3)
		returning id;
	`
	var id uint64
	err = r.db.QueryRowContext(ctx, query, dto.UserID, dto.Title, dataBytes).Scan(&id)
	if err != nil {
		r.logger.Log.Error("Ошибка при вставке секрета", zap.Uint64("user_id", dto.UserID), zap.String("title", dto.Title), zap.Error(err))
		return 0, fmt.Errorf("insert secret: %w", err)
	}

	r.logger.Log.Info("Секрет успешно создан", zap.Uint64("secret_id", id), zap.Uint64("user_id", dto.UserID))
	return id, nil
}

// GetByID возвращает секрет по его ID.
// Если секрет не найден — возвращает nil, nil.
func (r *SecretRepositoryImpl) GetByID(ctx context.Context, id uint64) (*models.ReadSecretDTO, error) {
	query := `
		select id, user_id, title, data, created_at, updated_at
		from secrets
		where id = $1;
	`

	var dto models.ReadSecretDTO
	var rawData []byte
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&dto.ID, &dto.UserID, &dto.Title, &rawData, &dto.CreatedAt, &dto.UpdatedAt)
	if err == sql.ErrNoRows {
		r.logger.Log.Warn("Секрет не найден по ID", zap.Uint64("secret_id", id))
		return nil, nil
	}
	if err != nil {
		r.logger.Log.Error("Ошибка при получении секрета по ID", zap.Uint64("secret_id", id), zap.Error(err))
		return nil, err
	}

	if err := json.Unmarshal(rawData, &dto.Data); err != nil {
		r.logger.Log.Error("Ошибка при анмаршалинге данных секрета", zap.Uint64("secret_id", id), zap.Error(err))
		return nil, ErrUnmarshalPayload
	}

	r.logger.Log.Info("Секрет успешно получен", zap.Uint64("secret_id", dto.ID))
	return &dto, nil
}

// GetAllByUser возвращает все секреты, принадлежащие пользователю.
func (r *SecretRepositoryImpl) GetAllByUser(ctx context.Context, userID uint64) ([]models.ReadSecretDTO, error) {
	query := `
		select id, user_id, title, data, created_at, updated_at
		from secrets
		where user_id = $1
		order by created_at desc;
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		r.logger.Log.Error("Ошибка при получении всех секретов пользователя", zap.Uint64("user_id", userID), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var secrets []models.ReadSecretDTO
	for rows.Next() {
		var dto models.ReadSecretDTO
		var rawData []byte

		if err := rows.Scan(&dto.ID, &dto.UserID, &dto.Title, &rawData, &dto.CreatedAt, &dto.UpdatedAt); err != nil {
			r.logger.Log.Error("Ошибка при чтении строки секрета", zap.Error(err))
			return nil, err
		}

		if err := json.Unmarshal(rawData, &dto.Data); err != nil {
			r.logger.Log.Error("Ошибка при анмаршалинге данных секрета", zap.Uint64("secret_id", dto.ID), zap.Error(err))
			return nil, ErrUnmarshalPayload
		}

		secrets = append(secrets, dto)
	}

	r.logger.Log.Info("Секреты пользователя успешно получены", zap.Uint64("user_id", userID), zap.Int("count", len(secrets)))
	return secrets, nil
}

// DeleteByID удаляет секрет по его ID.
func (r *SecretRepositoryImpl) DeleteByID(ctx context.Context, id uint64) error {
	query := `
		delete from secrets
		where id = $1;
	`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Log.Error("Ошибка при удалении секрета", zap.Uint64("secret_id", id), zap.Error(err))
		return err
	}

	r.logger.Log.Info("Секрет успешно удалён", zap.Uint64("secret_id", id))
	return nil
}
