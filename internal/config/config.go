package config

import (
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/shekshuev/gophkeeper/internal/logger"
	"go.uber.org/zap"
)

// Config содержит настройки приложения, включая параметры сервера, базы данных и токенов авторизации.
type Config struct {
	// ServerAddress — адрес и порт, на котором запускается сервер (например, "localhost:8080").
	ServerAddress string `env:"SERVER_ADDRESS"`

	// DatabaseDSN — строка подключения к базе данных (например, "host=localhost user=postgres dbname=gophkeeper sslmode=disable").
	DatabaseDSN string `env:"DATABASE_DSN"`

	// AccessTokenExpires — время жизни access токена (например, "15m", "1h").
	AccessTokenExpires time.Duration `env:"ACCESS_TOKEN_EXPIRES"`

	// RefreshTokenExpires — время жизни refresh токена.
	RefreshTokenExpires time.Duration `env:"REFRESH_TOKEN_EXPIRES"`

	// AccessTokenSecret — секретный ключ для подписи access токенов.
	AccessTokenSecret string `env:"ACCESS_TOKEN_SECRET"`

	// RefreshTokenSecret — секретный ключ для подписи refresh токенов.
	RefreshTokenSecret string `env:"REFRESH_TOKEN_SECRET"`
}

// GetConfig загружает конфигурацию из переменных окружения.
func GetConfig() Config {
	var cfg Config
	l := logger.NewLogger()
	err := env.Parse(&cfg)
	if err != nil {
		l.Log.Error("Error starting server", zap.Error(err))
	}
	return cfg
}
