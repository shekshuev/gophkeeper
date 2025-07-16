package models

import "time"

// CreateSecretDTO используется для создания нового секрета.
type CreateSecretDTO struct {
	UserID uint64        // ID владельца
	Title  string        // Название секрета
	Data   SecretDataDTO // Полезные данные (json)
}

// ReadSecretDTO используется для возврата секрета клиенту.
type ReadSecretDTO struct {
	ID        uint64        `json:"id"`         // ID секрета
	UserID    uint64        `json:"user_id"`    // ID владельца
	Title     string        `json:"title"`      // Название секрета
	Data      SecretDataDTO `json:"data"`       // Данные секрета
	CreatedAt time.Time     `json:"created_at"` // Когда создан
	UpdatedAt time.Time     `json:"updated_at"` // Когда обновлён
}

// SecretDataDTO представляет собой обёртку для различных типов приватных данных.
// Сохраняется как JSONB в базе.
type SecretDataDTO struct {
	LoginPassword *LoginPasswordData `json:"login_password,omitempty"` // Пара логин/пароль
	Text          *string            `json:"text,omitempty"`           // Произвольный текст
	Binary        []byte             `json:"binary,omitempty"`         // Бинарные данные
	Card          *CardData          `json:"card,omitempty"`           // Данные карты
}

// LoginPasswordData содержит логин и пароль.
type LoginPasswordData struct {
	Login    string `json:"login"`    // Логин
	Password string `json:"password"` // Пароль
}

// CardData содержит данные банковской карты.
type CardData struct {
	Number     string `json:"number"`      // Номер карты
	Holder     string `json:"holder"`      // Имя владельца
	ExpireDate string `json:"expire_date"` // Срок действия
	CVV        string `json:"cvv"`         // CVV-код
}
