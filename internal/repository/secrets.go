package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/stdlib"

	"github.com/shekshuev/gophkeeper/internal/config"
	"github.com/shekshuev/gophkeeper/internal/models"
)

// SecretRepositoryImpl реализует интерфейс SecretRepository с использованием PostgreSQL.
type SecretRepositoryImpl struct {
	db  *sql.DB
	cfg *config.Config
}

// NewSecretRepositoryImpl создаёт экземпляр репозитория секретов и устанавливает подключение к БД.
func NewSecretRepositoryImpl(cfg *config.Config) *SecretRepositoryImpl {
	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		log.Fatal("error connecting to database: ", err)
		return nil
	}
	return &SecretRepositoryImpl{cfg: cfg, db: db}
}

// Create сохраняет новый секрет в БД.
// Возвращает ID созданного секрета или ошибку.
func (r *SecretRepositoryImpl) Create(ctx context.Context, dto models.CreateSecretDTO) (uint64, error) {
	dataBytes, err := json.Marshal(dto.Data)
	if err != nil {
		return 0, ErrMarshalPayload
	}

	query := `
		insert into secrets (user_id, title, data)
		values ($1, $2, $3)
		returning id
	`
	var id uint64
	err = r.db.QueryRowContext(ctx, query, dto.UserID, dto.Title, dataBytes).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert secret: %w", err)
	}
	return id, nil
}

// GetByID возвращает секрет по его ID.
// Если секрет не найден — возвращается nil, nil.
func (r *SecretRepositoryImpl) GetByID(ctx context.Context, id uint64) (*models.ReadSecretDTO, error) {
	query := `
		select id, user_id, title, data, created_at, updated_at
		from secrets
		where id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var dto models.ReadSecretDTO
	var rawData []byte
	err := row.Scan(&dto.ID, &dto.UserID, &dto.Title, &rawData, &dto.CreatedAt, &dto.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	err = json.Unmarshal(rawData, &dto.Data)
	if err != nil {
		return nil, ErrUnmarshalPayload
	}
	return &dto, nil
}

// GetAllByUser возвращает все секреты, принадлежащие пользователю.
func (r *SecretRepositoryImpl) GetAllByUser(ctx context.Context, userID uint64) ([]models.ReadSecretDTO, error) {
	query := `
		select id, user_id, title, data, created_at, updated_at
		from secrets
		where user_id = $1
		order by created_at desc
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var secrets []models.ReadSecretDTO
	for rows.Next() {
		var dto models.ReadSecretDTO
		var rawData []byte
		err := rows.Scan(&dto.ID, &dto.UserID, &dto.Title, &rawData, &dto.CreatedAt, &dto.UpdatedAt)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(rawData, &dto.Data); err != nil {
			return nil, ErrUnmarshalPayload
		}
		secrets = append(secrets, dto)
	}
	return secrets, nil
}

// DeleteByID удаляет секрет по его ID.
func (r *SecretRepositoryImpl) DeleteByID(ctx context.Context, id uint64) error {
	query := `
		delete from secrets
		where id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
