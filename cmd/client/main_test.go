package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/shekshuev/gophkeeper/internal/client"
	"github.com/shekshuev/gophkeeper/internal/config"
	"github.com/shekshuev/gophkeeper/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestPrintBuildInfo(t *testing.T) {
	output := client.CaptureOutput(func() {
		printBuildInfo()
	})

	assert.Contains(t, output, "GophKeeper CLI")
	assert.Contains(t, output, "Версия сборки:")
	assert.Contains(t, output, "Дата сборки:")
	assert.Contains(t, output, "Коммит:")
}

func TestPrompt(t *testing.T) {
	r, w, _ := os.Pipe()
	origStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = origStdin }()

	go func() {
		fmt.Fprint(w, "  hello world  \n")
		w.Close()
	}()

	result := prompt("Введите: ")

	assert.Equal(t, "hello world", result)
}

func TestIsTokenValid_ValidToken(t *testing.T) {
	secret := "test-secret"
	tokenStr, err := utils.CreateToken(secret, "user-123", time.Minute)
	assert.NoError(t, err)

	ok := isTokenValid(
		func() (string, error) {
			return tokenStr, nil
		},
		func() config.Config {
			return config.Config{
				AccessTokenSecret: secret,
			}
		},
	)

	assert.True(t, ok)
}

func TestIsTokenValid_InvalidToken(t *testing.T) {
	ok := isTokenValid(
		func() (string, error) {
			return "bad.token", nil
		},
		func() config.Config {
			return config.Config{
				AccessTokenSecret: "test-secret",
			}
		},
	)

	assert.False(t, ok)
}

func TestIsTokenValid_EmptyToken(t *testing.T) {
	ok := isTokenValid(
		func() (string, error) {
			return "", nil
		},
		func() config.Config {
			return config.Config{}
		},
	)

	assert.False(t, ok)
}
