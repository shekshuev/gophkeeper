package client

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func captureOutput(f func()) string {
	var buf bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	buf.ReadFrom(r)
	os.Stdout = stdout
	return buf.String()
}

func TestGetSecret_Success(t *testing.T) {
	body := `{
		"id": 1,
		"user_id": 42,
		"title": "My Secret",
		"data": {},
		"created_at": "2025-07-16T00:00:00Z",
		"updated_at": "2025-07-16T00:00:00Z"
	}`

	client := resty.New()
	client.SetTransport(&mockRoundTripper{
		statusCode: 200,
		body:       body,
	})

	output := captureOutput(func() {
		GetSecret(1, client)
	})

	assert.Contains(t, output, `"title": "My Secret"`)
}

func TestGetSecret_NotFound(t *testing.T) {
	client := resty.New()
	client.SetTransport(&mockRoundTripper{
		statusCode: 404,
		body:       `Not found`,
	})

	output := captureOutput(func() {
		GetSecret(999, client)
	})

	assert.Contains(t, output, "Not found")
}

func TestListSecrets_Success(t *testing.T) {

	client := resty.New()
	client.SetTransport(&mockRoundTripper{
		statusCode: 200,
		body:       `[{"id":1,"title":"First"},{"id":2,"title":"Second"}]`,
	})

	output := captureOutput(func() {
		ListSecrets(client, func() (string, error) {
			return "42", nil
		})
	})

	assert.Contains(t, output, "1  First")
	assert.Contains(t, output, "2  Second")
}

func TestListSecrets_Unauthorized(t *testing.T) {

	output := captureOutput(func() {
		ListSecrets(resty.New(), func() (string, error) {
			return "", fmt.Errorf("unauthorized")
		})
	})

	assert.Contains(t, output, "Ошибка авторизации")
}

func TestDeleteSecret_Success(t *testing.T) {
	client := newMockClient(204, ``)

	output := captureOutput(func() {
		DeleteSecret(123, client)
	})

	assert.Contains(t, output, "204")
}

func TestDeleteSecret_Error(t *testing.T) {
	client := resty.New()
	client.SetTransport(&errorRoundTripper{})

	output := captureOutput(func() {
		DeleteSecret(123, client)
	})

	assert.Contains(t, output, "Ошибка:")
}

func TestCreateSecret_Text(t *testing.T) {
	restore := mockInput("1", "hello world")
	defer restore()

	client := newMockClient(201, `{"ok":true}`)

	output := captureOutput(func() {
		CreateSecret("MyText", client)
	})

	assert.Contains(t, output, "201")
}

func TestCreateSecret_LoginPassword(t *testing.T) {
	restore := mockInput("2", "admin", "pass123")
	defer restore()

	client := newMockClient(200, `{"ok":true}`)

	output := captureOutput(func() {
		CreateSecret("Creds", client)
	})

	assert.Contains(t, output, "200")
}

func TestCreateSecret_Card(t *testing.T) {
	restore := mockInput("3", "4111111111111111", "John Doe", "12/26", "123")
	defer restore()

	client := newMockClient(200, `{"ok":true}`)

	output := captureOutput(func() {
		CreateSecret("Card", client)
	})

	assert.Contains(t, output, "200")
}

func TestCreateSecret_Binary_NotSupported(t *testing.T) {
	restore := mockInput("4")
	defer restore()

	output := captureOutput(func() {
		CreateSecret("Binary", resty.New())
	})

	assert.Contains(t, output, "не реализованы")
}

func TestCreateSecret_InvalidType(t *testing.T) {
	restore := mockInput("999")
	defer restore()

	output := captureOutput(func() {
		CreateSecret("Unknown", resty.New())
	})

	assert.Contains(t, output, "Неверный тип")
}
