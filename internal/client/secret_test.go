package client

import (
	"fmt"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetSecret(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
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

		output := CaptureOutput(func() {
			GetSecret(1, client)
		})

		assert.Contains(t, output, `"title": "My Secret"`)
	})

	t.Run("NotFound", func(t *testing.T) {
		client := resty.New()
		client.SetTransport(&mockRoundTripper{
			statusCode: 404,
			body:       `Not found`,
		})

		output := CaptureOutput(func() {
			GetSecret(999, client)
		})

		assert.Contains(t, output, "Not found")
	})

	t.Run("HttpError", func(t *testing.T) {
		client := resty.New()
		client.SetTransport(&errorRoundTripper{})
		output := CaptureOutput(func() {
			GetSecret(1, client)
		})
		assert.Contains(t, output, "Ошибка:")
	})
}

func TestListSecrets(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		client := resty.New()
		client.SetTransport(&mockRoundTripper{
			statusCode: 200,
			body:       `[{"id":1,"title":"First"},{"id":2,"title":"Second"}]`,
		})

		output := CaptureOutput(func() {
			ListSecrets(client, func() (string, error) {
				return "42", nil
			})
		})

		assert.Contains(t, output, "1  First")
		assert.Contains(t, output, "2  Second")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		output := CaptureOutput(func() {
			ListSecrets(resty.New(), func() (string, error) {
				return "", fmt.Errorf("unauthorized")
			})
		})

		assert.Contains(t, output, "Ошибка авторизации")
	})

	t.Run("HttpError", func(t *testing.T) {
		client := resty.New()
		client.SetTransport(&errorRoundTripper{})
		output := CaptureOutput(func() {
			ListSecrets(client, func() (string, error) {
				return "42", nil
			})
		})
		assert.Contains(t, output, "Ошибка запроса")
	})
}

func TestDeleteSecret(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		client := newMockClient(204, ``)

		output := CaptureOutput(func() {
			DeleteSecret(123, client)
		})

		assert.Contains(t, output, "204")
	})

	t.Run("Error", func(t *testing.T) {
		client := resty.New()
		client.SetTransport(&errorRoundTripper{})

		output := CaptureOutput(func() {
			DeleteSecret(123, client)
		})

		assert.Contains(t, output, "Ошибка:")
	})
}

func TestCreateSecret(t *testing.T) {
	tcs := []struct {
		name   string
		inputs []string
		status int
		body   string
		expect string
	}{
		{
			name:   "Text",
			inputs: []string{"1", "hello world"},
			status: 201,
			body:   `{"ok":true}`,
			expect: "201",
		},
		{
			name:   "LoginPassword",
			inputs: []string{"2", "admin", "pass123"},
			status: 200,
			body:   `{"ok":true}`,
			expect: "200",
		},
		{
			name:   "Card",
			inputs: []string{"3", "4111111111111111", "John Doe", "12/26", "123"},
			status: 200,
			body:   `{"ok":true}`,
			expect: "200",
		},
		{
			name:   "Binary_NotSupported",
			inputs: []string{"4"},
			expect: "не реализованы",
		},
		{
			name:   "InvalidType",
			inputs: []string{"999"},
			expect: "Неверный тип",
		},
		{
			name:   "HttpError",
			inputs: []string{"1", "some text"},
			expect: "Ошибка запроса",
			status: -1,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			restore := MockInput(tc.inputs...)
			defer restore()

			var client *resty.Client
			if tc.status == -1 {
				client = resty.New()
				client.SetTransport(&errorRoundTripper{})
			} else {
				client = newMockClient(tc.status, tc.body)
			}

			output := CaptureOutput(func() {
				CreateSecret("SecretTitle", client)
			})

			assert.Contains(t, output, tc.expect)
		})
	}
}
