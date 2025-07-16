package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/shekshuev/gophkeeper/internal/models"
)

// CreateSecret — CLI-обёртка для создания нового секрета.
//
// Поддерживает типы:
//
//	[1] Произвольный текст
//	[2] Логин + пароль
//	[3] Банковская карта
//	[4] Бинарные данные (в hex или base64, пока не поддержано)
//
// Пользователь пошагово вводит данные через консоль, затем выполняется POST-запрос на /v1.0/secrets.
// В случае успеха выводится HTTP-статус и тело ответа.
func CreateSecret(title string, rc *resty.Client) {

	fmt.Println(`[1] Произвольный текст
[2] Логин + пароль
[3] Банковская карта
[4] Бинарные данные (в hex или base64, пока не поддержано)`)

	typ := prompt("Выберите тип секрета: ")

	var secretData models.SecretDataDTO

	switch typ {
	case "1":
		text := promptInput("Введите текст: ")
		secretData.Text = &text

	case "2":
		login := promptInput("Логин: ")
		password := promptInput("Пароль: ")
		secretData.LoginPassword = &models.LoginPasswordData{
			Login:    login,
			Password: password,
		}

	case "3":
		number := promptInput("Номер карты: ")
		holder := promptInput("Имя владельца: ")
		exp := promptInput("Срок действия (MM/YY): ")
		cvv := promptInput("CVV: ")
		secretData.Card = &models.CardData{
			Number:     number,
			Holder:     holder,
			ExpireDate: exp,
			CVV:        cvv,
		}

	case "4":
		fmt.Println("Бинарные данные пока не реализованы.")
		return

	default:
		fmt.Println("Неверный тип.")
		return
	}

	payload := models.CreateSecretDTO{
		Title: title,
		Data:  secretData,
	}

	resp, err := rc.R().
		SetBody(payload).
		Post("/v1.0/secrets")
	if err != nil {
		fmt.Println("Ошибка запроса:", err)
		return
	}
	fmt.Println(resp.StatusCode(), string(resp.Body()))
}

// GetSecret — CLI-обёртка для получения одного секрета по ID.
//
// Выполняет GET-запрос на /v1.0/secrets/{id}.
// В случае успеха выводит отформатированный JSON с данными секрета.
// В случае ошибки выводит статус и тело ответа.
func GetSecret(id uint64, rc *resty.Client) {

	var secret models.ReadSecretDTO
	resp, err := rc.R().
		SetResult(&secret).
		Get(fmt.Sprintf("/v1.0/secrets/%d", id))
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	if resp.IsError() {
		fmt.Println(resp.StatusCode(), string(resp.Body()))
		return
	}

	j, _ := json.MarshalIndent(secret, "", "  ")
	fmt.Println(string(j))
}

// ListSecrets — CLI-обёртка для получения всех секретов пользователя.
//
// Получает user_id с помощью функции getUserID (обычно из access-токена).
// Затем выполняет GET-запрос на /v1.0/secrets/user/{user_id}.
// В случае успеха выводит ID и название каждого секрета построчно.
func ListSecrets(rc *resty.Client, getUserID func() (string, error)) {

	userID, err := getUserID()
	if err != nil {
		fmt.Println("Ошибка авторизации:", err)
		return
	}

	var secrets []models.ReadSecretDTO
	resp, err := rc.R().
		SetResult(&secrets).
		Get("/v1.0/secrets/user/" + userID)
	if err != nil {
		fmt.Println("Ошибка запроса:", err)
		return
	}
	if resp.IsError() {
		fmt.Println(resp.StatusCode(), string(resp.Body()))
		return
	}

	for _, s := range secrets {
		fmt.Printf("%d  %s\n", s.ID, s.Title)
	}
}

// DeleteSecret — CLI-обёртка для удаления секрета по ID.
//
// Выполняет DELETE-запрос на /v1.0/secrets/{id}.
// В случае успеха выводит HTTP-статус. В случае ошибки — сообщение об ошибке.
func DeleteSecret(id uint64, rc *resty.Client) {

	resp, err := rc.R().
		Delete(fmt.Sprintf("/v1.0/secrets/%d", id))
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Println(resp.StatusCode())
}

// promptInput — вспомогательная функция для запроса строки от пользователя через консоль.
//
// Выводит переданную метку и считывает строку ввода.
// Возвращает trimmed-значение строки.
func promptInput(label string) string {
	fmt.Print(label)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}
