package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shekshuev/gophkeeper/internal/config"
	"github.com/shekshuev/gophkeeper/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestUserRepositoryImpl_CreateUser(t *testing.T) {
	testCases := []struct {
		name      string
		createDTO models.CreateUserDTO
		readDTO   models.ReadAuthUserDataDTO
		hasError  bool
	}{
		{
			name: "Success create",
			createDTO: models.CreateUserDTO{
				UserName:     "john",
				FirstName:    "John",
				LastName:     "Doe",
				PasswordHash: "password",
			},
			readDTO: models.ReadAuthUserDataDTO{
				ID:           1,
				UserName:     "john",
				PasswordHash: "password",
			},
			hasError: false,
		},
		{
			name: "Error on insert SQL",
			createDTO: models.CreateUserDTO{
				UserName:     "john",
				FirstName:    "John",
				LastName:     "Doe",
				PasswordHash: "password",
			},
			hasError: true,
		},
	}
	cfg := config.GetConfig()
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Error creating db mock: %v", err)
	}
	defer db.Close()
	r := &UserRepositoryImpl{cfg: &cfg, db: db}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.hasError {
				mock.ExpectQuery(regexp.QuoteMeta(`
					insert into users (user_name, first_name, last_name, password_hash) values ($1, $2, $3, $4) 
					returning id, user_name, password_hash;
					`)).
					WithArgs(
						tc.createDTO.UserName,
						tc.createDTO.FirstName,
						tc.createDTO.LastName,
						tc.createDTO.PasswordHash).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id", "user_name", "password_hash"},
						).AddRow(
							tc.readDTO.ID,
							tc.readDTO.UserName,
							tc.readDTO.PasswordHash,
						),
					)
			} else {
				mock.ExpectQuery(regexp.QuoteMeta(`
					insert into users (user_name, first_name, last_name, password_hash) values ($1, $2, $3, $4) 
					returning id, user_name, password_hash;
					`)).
					WithArgs(
						tc.createDTO.UserName,
						tc.createDTO.FirstName,
						tc.createDTO.LastName,
						tc.createDTO.PasswordHash).
					WillReturnError(sql.ErrNoRows)
			}
			ctx := context.Background()
			user, err := r.CreateUser(ctx, tc.createDTO)
			if tc.hasError {
				assert.NotNil(t, err, "Error is nil")
			} else {
				assert.Nil(t, err, "Error is not nil")
				assert.Equal(t, tc.readDTO, *user, "User mismatch")
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Not all expectations were met: %v", err)
			}
		})
	}
}

func TestUserRepositoryImpl_GetUserByUserName(t *testing.T) {
	testCases := []struct {
		name     string
		userName string
		readDTO  *models.ReadAuthUserDataDTO
		hasError bool
	}{
		{
			name:     "Success get user by ID",
			userName: "john",
			readDTO: &models.ReadAuthUserDataDTO{
				ID:           1,
				UserName:     "john",
				PasswordHash: "password",
			},
			hasError: false,
		},
		{
			name:     "User not found",
			userName: "notfound",
			readDTO:  nil,
			hasError: true,
		},
	}

	cfg := config.GetConfig()
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Error creating db mock: %v", err)
	}
	defer db.Close()

	r := &UserRepositoryImpl{cfg: &cfg, db: db}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.hasError {
				rows := sqlmock.NewRows([]string{
					"id", "user_name", "password_hash",
				}).AddRow(
					tc.readDTO.ID, tc.readDTO.UserName, tc.readDTO.PasswordHash,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`select id, user_name, password_hash from users where user_name = $1 and deleted_at is null;`)).
					WithArgs(tc.userName).
					WillReturnRows(rows)
			} else {
				mock.ExpectQuery(regexp.QuoteMeta(`select id, user_name, password_hash from users where user_name = $1 and deleted_at is null;`)).
					WithArgs(tc.userName).
					WillReturnError(sql.ErrNoRows)
			}

			ctx := context.Background()
			user, err := r.GetUserByUserName(ctx, tc.userName)
			if tc.hasError {
				assert.NotNil(t, err, "Error is nil")
				assert.Nil(t, user, "User should be nil")
			} else {
				assert.Nil(t, err, "Error is not nil")
				assert.Equal(t, tc.readDTO, user, "User mismatch")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Not all expectations were met: %v", err)
			}
		})
	}
}

func TestUserRepositoryImpl_GetUserByID(t *testing.T) {
	testCases := []struct {
		name     string
		id       uint64
		readDTO  *models.ReadUserDTO
		hasError bool
	}{
		{
			name: "Success get user by ID",
			id:   1,
			readDTO: &models.ReadUserDTO{
				ID:        1,
				UserName:  "john",
				FirstName: "John",
				LastName:  "Doe",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			hasError: false,
		},
		{
			name:     "User not found",
			id:       999,
			readDTO:  nil,
			hasError: true,
		},
	}

	cfg := config.GetConfig()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating db mock: %v", err)
	}
	defer db.Close()

	r := &UserRepositoryImpl{cfg: &cfg, db: db}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			if !tc.hasError {
				mock.ExpectQuery(regexp.QuoteMeta(`
					select 
						id, user_name, first_name, last_name, created_at, updated_at 
					from users 
					where id = $1 and deleted_at is null;
				`)).
					WithArgs(tc.id).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "user_name", "first_name", "last_name", "created_at", "updated_at",
					}).AddRow(
						tc.readDTO.ID,
						tc.readDTO.UserName,
						tc.readDTO.FirstName,
						tc.readDTO.LastName,
						tc.readDTO.CreatedAt,
						tc.readDTO.UpdatedAt,
					))
			} else {
				mock.ExpectQuery(regexp.QuoteMeta(`
					select 
						id, user_name, first_name, last_name, created_at, updated_at 
					from users 
					where id = $1 and deleted_at is null;
				`)).
					WithArgs(tc.id).
					WillReturnError(sql.ErrNoRows)
			}

			user, err := r.GetUserByID(ctx, tc.id)

			if tc.hasError {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.readDTO.ID, user.ID)
				assert.Equal(t, tc.readDTO.UserName, user.UserName)
				assert.Equal(t, tc.readDTO.FirstName, user.FirstName)
				assert.Equal(t, tc.readDTO.LastName, user.LastName)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
