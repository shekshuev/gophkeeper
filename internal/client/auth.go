package client

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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
// Выполняет POST-запрос к API /v1.0/auth/register.
// Выводит сообщение об успехе или ошибке.
func Register() {
	user := models.RegisterUserDTO{
		UserName:        prompt("Имя пользователя: "),
		Password:        prompt("Пароль: "),
		PasswordConfirm: prompt("Подтвердите пароль: "),
		FirstName:       prompt("Имя: "),
		LastName:        prompt("Фамилия: "),
	}

	rc := api()

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
// Запрашивает у пользователя логин и пароль,
// отправляет их на эндпоинт /v1.0/auth/login.
//
// При успешной авторизации сохраняет access-токен локально и выводит сообщение об успешном входе.
func Login() {
	user := models.LoginUserDTO{
		UserName: prompt("Username: "),
		Password: prompt("Password: "),
	}

	rc := api()

	var tokens models.ReadTokenDTO

	resp, err := rc.R().
		SetResult(&tokens).
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

	if err := SaveToken(tokens.AccessToken); err != nil {
		fmt.Println("Не удалось сохранить токен:", err)
		return
	}

	fmt.Println("Вход выполнен.")
}
