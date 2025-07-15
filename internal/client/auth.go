package client

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/shekshuev/gophkeeper/internal/models"
)

func prompt(label string) string {
	fmt.Print(label)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

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
