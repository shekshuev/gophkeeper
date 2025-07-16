package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shekshuev/gophkeeper/internal/config"
	"github.com/shekshuev/gophkeeper/internal/logger"
	"github.com/shekshuev/gophkeeper/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestSecretRepositoryImpl_Create(t *testing.T) {
	cfg := config.GetConfig()
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &SecretRepositoryImpl{cfg: &cfg, db: db, logger: logger.NewLogger()}

	dto := models.CreateSecretDTO{
		UserID: 42,
		Title:  "Login creds",
		Data: models.SecretDataDTO{
			LoginPassword: &models.LoginPasswordData{
				Login:    "admin",
				Password: "1234",
			},
		},
	}

	dataBytes, _ := json.Marshal(dto.Data)

	mock.ExpectQuery(regexp.QuoteMeta(`
		insert into secrets (user_id, title, data)
		values ($1, $2, $3)
		returning id
	`)).
		WithArgs(dto.UserID, dto.Title, dataBytes).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uint64(1)))

	id, err := repo.Create(context.Background(), dto)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSecretRepositoryImpl_Create_QueryError(t *testing.T) {
	cfg := config.GetConfig()
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dto := models.CreateSecretDTO{
		UserID: 1,
		Title:  "Test",
		Data:   models.SecretDataDTO{Text: ptr("secret")},
	}
	dataBytes, _ := json.Marshal(dto.Data)

	mock.ExpectQuery("insert into secrets").
		WithArgs(dto.UserID, dto.Title, dataBytes).
		WillReturnError(assert.AnError)

	repo := &SecretRepositoryImpl{cfg: &cfg, db: db, logger: logger.NewLogger()}
	id, err := repo.Create(context.Background(), dto)
	assert.Error(t, err)
	assert.Equal(t, uint64(0), id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSecretRepositoryImpl_GetByID(t *testing.T) {
	cfg := config.GetConfig()
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &SecretRepositoryImpl{cfg: &cfg, db: db, logger: logger.NewLogger()}

	now := time.Now()
	rawData := models.SecretDataDTO{
		Text: ptr("some secret text"),
	}
	dataBytes, _ := json.Marshal(rawData)

	mock.ExpectQuery(regexp.QuoteMeta(`
		select id, user_id, title, data, created_at, updated_at
		from secrets
		where id = $1
	`)).
		WithArgs(uint64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "title", "data", "created_at", "updated_at",
		}).AddRow(
			uint64(1), uint64(42), "Note", dataBytes, now, now,
		))

	secret, err := repo.GetByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, secret)
	assert.Equal(t, uint64(1), secret.ID)
	assert.Equal(t, "Note", secret.Title)
	assert.NotNil(t, secret.Data.Text)
	assert.Equal(t, "some secret text", *secret.Data.Text)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSecretRepositoryImpl_GetByID_NotFound(t *testing.T) {
	cfg := config.GetConfig()
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("select id, user_id, title").
		WithArgs(uint64(999)).
		WillReturnError(sql.ErrNoRows)

	repo := &SecretRepositoryImpl{cfg: &cfg, db: db, logger: logger.NewLogger()}
	secret, err := repo.GetByID(context.Background(), 999)
	assert.NoError(t, err)
	assert.Nil(t, secret)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSecretRepositoryImpl_GetByID_QueryError(t *testing.T) {
	cfg := config.GetConfig()
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("select id, user_id, title").
		WithArgs(uint64(1)).
		WillReturnError(assert.AnError)

	repo := &SecretRepositoryImpl{cfg: &cfg, db: db, logger: logger.NewLogger()}
	secret, err := repo.GetByID(context.Background(), 1)
	assert.Error(t, err)
	assert.Nil(t, secret)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSecretRepositoryImpl_GetByID_UnmarshalError(t *testing.T) {
	cfg := config.GetConfig()
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	now := time.Now()
	mock.ExpectQuery("select id, user_id, title").
		WithArgs(uint64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "title", "data", "created_at", "updated_at",
		}).AddRow(1, 42, "Broken", []byte("not-json"), now, now))

	repo := &SecretRepositoryImpl{cfg: &cfg, db: db, logger: logger.NewLogger()}
	secret, err := repo.GetByID(context.Background(), 1)
	assert.ErrorIs(t, err, ErrUnmarshalPayload)
	assert.Nil(t, secret)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSecretRepositoryImpl_GetAllByUser(t *testing.T) {
	cfg := config.GetConfig()
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &SecretRepositoryImpl{cfg: &cfg, db: db, logger: logger.NewLogger()}

	now := time.Now()
	data := models.SecretDataDTO{
		Card: &models.CardData{
			Number:     "1234567890123456",
			Holder:     "John Doe",
			ExpireDate: "12/25",
			CVV:        "123",
		},
	}
	dataBytes, _ := json.Marshal(data)

	mock.ExpectQuery(regexp.QuoteMeta(`
		select id, user_id, title, data, created_at, updated_at
		from secrets
		where user_id = $1
		order by created_at desc
	`)).
		WithArgs(uint64(42)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "title", "data", "created_at", "updated_at",
		}).
			AddRow(uint64(1), uint64(42), "Card 1", dataBytes, now, now).
			AddRow(uint64(2), uint64(42), "Card 2", dataBytes, now, now),
		)

	secrets, err := repo.GetAllByUser(context.Background(), 42)
	assert.NoError(t, err)
	assert.Len(t, secrets, 2)
	assert.Equal(t, "Card 1", secrets[0].Title)
	assert.NotNil(t, secrets[0].Data.Card)
	assert.Equal(t, "John Doe", secrets[0].Data.Card.Holder)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSecretRepositoryImpl_GetAllByUser_QueryError(t *testing.T) {
	cfg := config.GetConfig()
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("select id, user_id, title").
		WithArgs(uint64(1)).
		WillReturnError(assert.AnError)

	repo := &SecretRepositoryImpl{cfg: &cfg, db: db, logger: logger.NewLogger()}
	secrets, err := repo.GetAllByUser(context.Background(), 1)
	assert.Error(t, err)
	assert.Nil(t, secrets)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSecretRepositoryImpl_GetAllByUser_ScanError(t *testing.T) {
	cfg := config.GetConfig()
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("select id, user_id, title").
		WithArgs(uint64(42)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id",
		}).AddRow(1, 42))

	repo := &SecretRepositoryImpl{cfg: &cfg, db: db, logger: logger.NewLogger()}
	secrets, err := repo.GetAllByUser(context.Background(), 42)
	assert.Error(t, err)
	assert.Nil(t, secrets)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSecretRepositoryImpl_GetAllByUser_UnmarshalError(t *testing.T) {
	cfg := config.GetConfig()
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	now := time.Now()

	mock.ExpectQuery("select id, user_id, title").
		WithArgs(uint64(42)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "title", "data", "created_at", "updated_at",
		}).AddRow(1, 42, "Bad", []byte("broken-json"), now, now))

	repo := &SecretRepositoryImpl{cfg: &cfg, db: db, logger: logger.NewLogger()}
	secrets, err := repo.GetAllByUser(context.Background(), 42)
	assert.ErrorIs(t, err, ErrUnmarshalPayload)
	assert.Nil(t, secrets)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSecretRepositoryImpl_DeleteByID(t *testing.T) {
	cfg := config.GetConfig()
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &SecretRepositoryImpl{cfg: &cfg, db: db, logger: logger.NewLogger()}

	mock.ExpectExec(regexp.QuoteMeta(`
		delete from secrets
		where id = $1
	`)).
		WithArgs(uint64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.DeleteByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSecretRepositoryImpl_DeleteByID_Error(t *testing.T) {
	cfg := config.GetConfig()
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec("delete from secrets").
		WithArgs(uint64(123)).
		WillReturnError(assert.AnError)

	repo := &SecretRepositoryImpl{cfg: &cfg, db: db, logger: logger.NewLogger()}
	err = repo.DeleteByID(context.Background(), 123)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func ptr(s string) *string {
	return &s
}
