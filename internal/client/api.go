package client

import (
	"github.com/go-resty/resty/v2"
	"github.com/shekshuev/gophkeeper/internal/config"
)

// api возвращает сконфигурированный resty.Client.
func api() *resty.Client {
	cfg := config.GetConfig()

	rc := resty.New().
		SetBaseURL("http://" + cfg.ServerAddress)

	token, _ := LoadToken()
	if token != "" {
		rc.SetHeader("Authorization", "Bearer "+token)
	}
	return rc
}
