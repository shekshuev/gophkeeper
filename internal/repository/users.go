package repository

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/stdlib"

	"github.com/shekshuev/gophkeeper/internal/config"
	"github.com/shekshuev/gophkeeper/internal/models"
)

// UserRepositoryImpl — реализация интерфейса UserRepository для работы с пользователями через SQL.
type UserRepositoryImpl struct {
	db  *sql.DB        // соединение с базой данных
	cfg *config.Config // конфигурация приложения
}

// NewUserRepositoryImpl создаёт новый экземпляр UserRepositoryImpl.
// Устанавливает соединение с базой данных на основе переданного DSN из конфигурации.
func NewUserRepositoryImpl(cfg *config.Config) *UserRepositoryImpl {
	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		log.Fatal("Error connecting to database", err)
		return nil
	}
	repository := &UserRepositoryImpl{cfg: cfg, db: db}
	return repository
}

// CreateUser добавляет нового пользователя в базу данных.
// Принимает DTO с необходимыми полями, возвращает DTO с ID, userName и хешем пароля.
// В случае ошибки возвращает её.
func (r *UserRepositoryImpl) CreateUser(ctx context.Context, dto models.CreateUserDTO) (*models.ReadAuthUserDataDTO, error) {
	query := `
		insert into users (user_name, first_name, last_name, password_hash) values ($1, $2, $3, $4)
		returning id, user_name, password_hash;
	`
	var user models.ReadAuthUserDataDTO
	err := r.db.QueryRowContext(
		ctx, query, dto.UserName, dto.FirstName, dto.LastName, dto.PasswordHash).Scan(&user.ID, &user.UserName, &user.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUserName получает пользователя по его userName, если он не помечен как удалённый.
// Возвращает ReadAuthUserDataDTO или ErrNotFound, если пользователь не найден.
func (r *UserRepositoryImpl) GetUserByUserName(ctx context.Context, userName string) (*models.ReadAuthUserDataDTO, error) {
	query := `
		select 
			id, user_name, password_hash 
		from users where user_name = $1 and deleted_at is null;
	`
	var user models.ReadAuthUserDataDTO
	err := r.db.QueryRowContext(ctx, query, userName).Scan(&user.ID, &user.UserName, &user.PasswordHash)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID получает пользователя по его уникальному идентификатору, если он не помечен как удалённый.
// Возвращает ReadUserDTO или ErrNotFound, если пользователь не найден.
func (r *UserRepositoryImpl) GetUserByID(ctx context.Context, id uint64) (*models.ReadUserDTO, error) {
	query := `
		select 
			id, user_name, first_name, last_name, created_at, updated_at 
		from users 
		where id = $1 and deleted_at is null;
	`
	var user models.ReadUserDTO
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&user.ID, &user.UserName, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
