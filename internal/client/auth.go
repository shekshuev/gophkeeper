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

// prompt отображает текстовый вопрос пользователю и считывает строку с консоли.
func prompt(label string) string {
	fmt.Print(label)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

// Register — CLI-обёртка для регистрации нового пользователя.
//
// Пошагово запрашивает у пользователя:
//   - имя пользователя
//   - пароль и его подтверждение
//   - имя и фамилию
//
// Отправляет POST-запрос к API /v1.0/auth/register через переданный HTTP-клиент.
//
// Аргументы:
//   - rc: настроенный HTTP-клиент (resty.Client)
//
// При ошибке запроса или ответа выводит сообщение в консоль.
// В случае успеха выводит сообщение "Успешно зарегистрирован."
func Register(rc *resty.Client) {
	user := models.RegisterUserDTO{
		UserName:        prompt("Имя пользователя: "),
		Password:        prompt("Пароль: "),
		PasswordConfirm: prompt("Подтвердите пароль: "),
		FirstName:       prompt("Имя: "),
		LastName:        prompt("Фамилия: "),
	}

	resp, err := rc.R().
		SetBody(user).
		Post("/v1.0/auth/register")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	if resp.IsError() {
		fmt.Println("Ошибка:", resp.Status(), string(resp.Body()))
		return
	}
	fmt.Println("Успешно зарегистрирован.")
}

// Login — CLI-обёртка для авторизации пользователя.
//
// Пошагово запрашивает у пользователя:
//   - имя пользователя
//   - пароль
//
// Отправляет POST-запрос к API /v1.0/auth/login через переданный HTTP-клиент.
//
// Аргументы:
//   - rc: настроенный HTTP-клиент (resty.Client)
//   - saveTokenFunc: функция сохранения access-токена (например, запись в файл)
//
// В случае ошибки (сетевой, HTTP или разбора JSON) выводит сообщение об ошибке.
// При успешной авторизации вызывает saveTokenFunc и выводит "Вход выполнен."
func Login(rc *resty.Client, saveTokenFunc func(string) error) {
	user := models.LoginUserDTO{
		UserName: prompt("Username: "),
		Password: prompt("Password: "),
	}

	var tokens models.ReadTokenDTO

	resp, err := rc.R().
		SetBody(user).
		Post("/v1.0/auth/login")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	if resp.IsError() {
		fmt.Println("Ошибка:", resp.Status(), string(resp.Body()))
		return
	}

	if err := json.Unmarshal(resp.Body(), &tokens); err != nil {
		fmt.Println("Ошибка разбора ответа:", err)
		return
	}

	if err := saveTokenFunc(tokens.AccessToken); err != nil {
		fmt.Println("Не удалось сохранить токен:", err)
		return
	}

	fmt.Println("Вход выполнен.")
}
