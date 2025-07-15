package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/shekshuev/gophkeeper/internal/models"
)

func CreateSecret(title string) {
	rc := api()

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
	fmt.Println(resp.Status(), string(resp.Body()))
}

func GetSecret(id uint64) {
	rc := api()

	var secret models.ReadSecretDTO
	resp, err := rc.R().
		SetResult(&secret).
		Get(fmt.Sprintf("/v1.0/secrets/%d", id))
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	if resp.IsError() {
		fmt.Println(resp.Status(), string(resp.Body()))
		return
	}

	j, _ := json.MarshalIndent(secret, "", "  ")
	fmt.Println(string(j))
}

func ListSecrets() {
	rc := api()

	userID, err := GetUserIDFromToken()
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
		fmt.Println(resp.Status(), string(resp.Body()))
		return
	}

	for _, s := range secrets {
		fmt.Printf("%d  %s\n", s.ID, s.Title)
	}
}

func DeleteSecret(id uint64) {
	rc := api()

	resp, err := rc.R().
		Delete(fmt.Sprintf("/v1.0/secrets/%d", id))
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Println(resp.Status())
}

// promptInput — простой ввод из терминала
func promptInput(label string) string {
	fmt.Print(label)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}
