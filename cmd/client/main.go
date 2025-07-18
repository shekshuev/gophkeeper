package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "github.com/joho/godotenv/autoload"

	"github.com/shekshuev/gophkeeper/internal/client"
	"github.com/shekshuev/gophkeeper/internal/config"
	"github.com/shekshuev/gophkeeper/internal/utils"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func printBuildInfo() {
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}

	fmt.Println("GophKeeper CLI")
	fmt.Printf("Версия сборки: %s\n", buildVersion)
	fmt.Printf("Дата сборки: %s\n", buildDate)
	fmt.Printf("Коммит: %s\n\n", buildCommit)
}

func prompt(label string) string {
	fmt.Print(label)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func isTokenValid(loadToken func() (string, error), getConfig func() config.Config) bool {
	cfg := getConfig()
	tokenStr, err := loadToken()
	if err != nil || tokenStr == "" {
		return false
	}

	claims, err := utils.GetToken(tokenStr, cfg.AccessTokenSecret)
	if err != nil || claims.Subject == "" {
		return false
	}
	return true
}

func isTokenValidDefault() bool {
	return isTokenValid(client.LoadToken, config.GetConfig)
}

func mainMenu() bool {
	for {
		fmt.Println(`[1] Создать секрет
[2] Показать все секреты
[3] Получить секрет по ID
[4] Удалить секрет по ID
[5] Завершить сессию
[0] Выйти`)
		choice := prompt("Выберите действие > ")

		switch choice {
		case "1":
			title := prompt("Введите название секрета: ")
			client.CreateSecret(title, client.Api())
		case "2":
			client.ListSecrets(client.Api(), client.GetUserIDFromToken)
		case "3":
			idStr := prompt("Введите ID секрета: ")
			id, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil {
				fmt.Println("Некорректный ID")
				continue
			}
			client.GetSecret(id, client.Api())
		case "4":
			idStr := prompt("Введите ID секрета для удаления: ")
			id, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil {
				fmt.Println("Некорректный ID")
				continue
			}
			client.DeleteSecret(id, client.Api())
		case "5":
			err := client.Logout()
			if err != nil {
				fmt.Println("Не удалось удалить токен:", err)
			} else {
				fmt.Println("Сессия завершена.")
			}
			return false
		case "0":
			fmt.Println("До свидания!")
			os.Exit(0)
		default:
			fmt.Println("Неизвестная команда")
		}
		fmt.Println()
	}
}

func authMenu() bool {
	for {
		fmt.Println(`[1] Зарегистрироваться
[2] Войти
[0] Выйти`)
		choice := prompt("Выберите действие > ")

		switch choice {
		case "1":
			client.Register(client.Api())
		case "2":
			client.Login(client.Api(), client.SaveToken)
			if isTokenValidDefault() {
				fmt.Println("Вход выполнен успешно.")
				return true
			}
		case "0":
			fmt.Println("До свидания!")
			os.Exit(0)
		default:
			fmt.Println("Неизвестная команда")
		}
		fmt.Println()
	}
}

func main() {
	printBuildInfo()

	for {
		if isTokenValidDefault() {
			if !mainMenu() {
				continue
			}
		} else {
			if authMenu() {
				continue
			}
		}
	}
}
