package client

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPI_BaseURL(t *testing.T) {
	os.Setenv("SERVER_ADDRESS", "localhost:9999")

	client := api()

	assert.Equal(t, "http://localhost:9999", client.BaseURL)
}

func TestAPI_WithToken(t *testing.T) {
	home, _ := os.UserHomeDir()
	tokenFile := filepath.Join(home, ".gophkeeper", "token.json")
	_ = os.MkdirAll(filepath.Dir(tokenFile), 0700)
	defer os.Remove(tokenFile)

	_ = os.WriteFile(tokenFile, []byte(`{"token": "test-token-123"}`), 0600)

	client := api()

	authHeader := client.Header.Get("Authorization")
	assert.Equal(t, "Bearer test-token-123", authHeader)
}
