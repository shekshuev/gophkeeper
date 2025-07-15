package client

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/shekshuev/gophkeeper/internal/config"
	"github.com/shekshuev/gophkeeper/internal/utils"
	"github.com/stretchr/testify/assert"
)

const testToken = "test.jwt.token"

func TestSaveAndLoadToken(t *testing.T) {
	_ = os.Remove(tokenPath())

	err := SaveToken(testToken)
	assert.NoError(t, err)

	loaded, err := LoadToken()
	assert.NoError(t, err)
	assert.Equal(t, testToken, loaded)

	err = Logout()
	assert.NoError(t, err)

	_, err = os.Stat(tokenPath())
	assert.True(t, os.IsNotExist(err))
}

func TestLoadToken_FileMissing(t *testing.T) {
	_ = os.Remove(tokenPath())
	_, err := LoadToken()
	assert.Error(t, err)
}

func TestLoadToken_InvalidJSON(t *testing.T) {
	_ = os.MkdirAll(filepath.Dir(tokenPath()), 0700)
	_ = os.WriteFile(tokenPath(), []byte("invalid json"), 0600)

	_, err := LoadToken()
	assert.Error(t, err)
}

func TestGetUserIDFromToken_Success(t *testing.T) {
	t.Setenv("ACCESS_TOKEN_SECRET", "mysecret")
	cfg := config.GetConfig()

	tokenStr, err := utils.CreateToken(cfg.AccessTokenSecret, "42", time.Hour)
	assert.NoError(t, err)

	err = SaveToken(tokenStr)
	assert.NoError(t, err)

	userID, err := GetUserIDFromToken()
	assert.NoError(t, err)
	assert.Equal(t, "42", userID)
}

func TestGetUserIDFromToken_InvalidToken(t *testing.T) {
	t.Setenv("ACCESS_TOKEN_SECRET", "mysecret")
	_ = SaveToken("not.a.real.jwt")

	_, err := GetUserIDFromToken()
	assert.Error(t, err)
}

func TestLogout_FileDoesNotExist(t *testing.T) {
	_ = os.Remove(tokenPath())
	err := Logout()
	assert.Error(t, err)
}
