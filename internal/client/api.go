package client

import (
	"github.com/go-resty/resty/v2"
	"github.com/shekshuev/gophkeeper/internal/config"
)

// Api возвращает сконфигурированный HTTP-клиент resty.Client с базовым URL и заголовком авторизации.
//
// Клиент автоматически:
//   - Устанавливает базовый адрес сервера из конфигурации (`cfg.ServerAddress`).
//   - Добавляет заголовок `Authorization: Bearer <токен>`, если токен ранее сохранён.
//
// Используется везде, где требуется выполнять HTTP-запросы к API сервера.
func Api() *resty.Client {
	cfg := config.GetConfig()

	rc := resty.New().
		SetBaseURL("http://" + cfg.ServerAddress)

	token, _ := LoadToken()
	if token != "" {
		rc.SetHeader("Authorization", "Bearer "+token)
	}

	return rc
}
