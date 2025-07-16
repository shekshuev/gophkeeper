package client

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/shekshuev/gophkeeper/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestRegister_Success(t *testing.T) {
	defer mockInput("testuser", "secret", "secret", "John", "Doe")()

	client := newMockClient(200, `{}`)

	var out bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Register(client)

	w.Close()
	out.ReadFrom(r)
	os.Stdout = stdout

	assert.Contains(t, out.String(), "Успешно зарегистрирован.")
}

func TestRegister_Fail(t *testing.T) {
	defer mockInput("testuser", "secret", "secret", "John", "Doe")()

	client := newMockClient(400, `{"error":"Username already exists"}`)

	var out bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Register(client)

	w.Close()
	out.ReadFrom(r)
	os.Stdout = stdout

	assert.Contains(t, out.String(), "Ошибка:")
}

func TestRegister_HttpError(t *testing.T) {
	defer mockInput("u", "p", "p", "f", "l")()

	client := resty.New()
	client.SetTransport(&mockRoundTripper{
		statusCode: 500,
		body:       `internal error`,
	})

	var out bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Register(client)

	w.Close()
	out.ReadFrom(r)
	os.Stdout = stdout

	assert.Contains(t, out.String(), "Ошибка:")
}

func TestLogin_Success(t *testing.T) {
	defer mockInput("testuser", "secret")()

	tokens := &models.ReadTokenDTO{}
	client := resty.New()
	client.SetTransport(&mockRoundTripper{
		statusCode: 200,
		body:       `{"access_token":"mocked-token"}`,
	})
	client.OnAfterResponse(func(c *resty.Client, res *resty.Response) error {
		return json.Unmarshal(res.Body(), tokens)
	})

	var out bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Login(client, func(token string) error {
		assert.Equal(t, "mocked-token", token)
		return nil
	})

	w.Close()
	out.ReadFrom(r)
	os.Stdout = stdout

	assert.Contains(t, out.String(), "Вход выполнен.")
}

func TestLogin_Invalid(t *testing.T) {
	defer mockInput("testuser", "wrongpass")()

	client := newMockClient(401, `Unauthorized`)

	var out bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Login(client, func(token string) error {
		return nil
	})

	w.Close()
	out.ReadFrom(r)
	os.Stdout = stdout

	assert.Contains(t, out.String(), "Ошибка:")
}

func TestLogin_UnmarshalError(t *testing.T) {
	defer mockInput("u", "p")()

	client := resty.New()
	client.SetTransport(&mockRoundTripper{
		statusCode: 200,
		body:       `not-json`,
	})

	var out bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Login(client, func(token string) error { return nil })

	w.Close()
	out.ReadFrom(r)
	os.Stdout = stdout

	assert.Contains(t, out.String(), "Ошибка разбора ответа")
}

func TestLogin_RequestError(t *testing.T) {
	defer mockInput("u", "p")()

	client := resty.New()
	client.SetTransport(&errorRoundTripper{})

	var out bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Login(client, func(token string) error { return nil })

	w.Close()
	out.ReadFrom(r)
	os.Stdout = stdout

	assert.Contains(t, out.String(), "error:")
}
