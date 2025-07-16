package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	serverAddress := "localhost:3000"
	databaseDSN := "host=test port=5432 user=test password=test dbname=test sslmode=disable"
	accessTokenExpires := "15m"
	refreshTokenExpires := "720h"
	accessTokenSecret := "secret123"
	refreshTokenSecret := "refreshsecret456"

	os.Setenv("SERVER_ADDRESS", serverAddress)
	os.Setenv("DATABASE_DSN", databaseDSN)
	os.Setenv("ACCESS_TOKEN_EXPIRES", accessTokenExpires)
	os.Setenv("REFRESH_TOKEN_EXPIRES", refreshTokenExpires)
	os.Setenv("ACCESS_TOKEN_SECRET", accessTokenSecret)
	os.Setenv("REFRESH_TOKEN_SECRET", refreshTokenSecret)

	defer func() {
		os.Unsetenv("SERVER_ADDRESS")
		os.Unsetenv("DATABASE_DSN")
		os.Unsetenv("ACCESS_TOKEN_EXPIRES")
		os.Unsetenv("REFRESH_TOKEN_EXPIRES")
		os.Unsetenv("ACCESS_TOKEN_SECRET")
		os.Unsetenv("REFRESH_TOKEN_SECRET")
	}()

	cfg := GetConfig()
	assert.Equal(t, serverAddress, cfg.ServerAddress)
	assert.Equal(t, databaseDSN, cfg.DatabaseDSN)
	assert.Equal(t, 15*time.Minute, cfg.AccessTokenExpires)
	assert.Equal(t, 30*24*time.Hour, cfg.RefreshTokenExpires)
	assert.Equal(t, accessTokenSecret, cfg.AccessTokenSecret)
	assert.Equal(t, refreshTokenSecret, cfg.RefreshTokenSecret)
}
