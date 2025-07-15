package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Username string `validate:"alphanumunderscore,startswithalpha"`
	Password string `validate:"password"`
}

func TestValidator_ValidData(t *testing.T) {
	v := NewValidator()

	data := testStruct{
		Username: "John_Doe123",
		Password: "Pa$$w0rd!",
	}

	err := v.Struct(data)
	assert.NoError(t, err)
}

func TestValidator_Invalid_AlphanumUnderscore(t *testing.T) {
	v := NewValidator()

	data := testStruct{
		Username: "John-Doe!",
		Password: "Pa$$w0rd!",
	}

	err := v.Struct(data)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Username")
}

func TestValidator_Invalid_StartsWithAlpha(t *testing.T) {
	v := NewValidator()

	data := testStruct{
		Username: "123User",
		Password: "Pa$$w0rd!",
	}

	err := v.Struct(data)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Username")
}

func TestValidator_Invalid_Password_MissingDigit(t *testing.T) {
	v := NewValidator()

	data := testStruct{
		Username: "ValidUser",
		Password: "Password!",
	}

	err := v.Struct(data)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Password")
}

func TestValidator_Invalid_Password_MissingLetter(t *testing.T) {
	v := NewValidator()

	data := testStruct{
		Username: "ValidUser",
		Password: "1234567!",
	}

	err := v.Struct(data)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Password")
}

func TestValidator_Invalid_Password_MissingSpecial(t *testing.T) {
	v := NewValidator()

	data := testStruct{
		Username: "ValidUser",
		Password: "Passw0rd",
	}

	err := v.Struct(data)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Password")
}

func TestValidator_Invalid_Password_TooShort(t *testing.T) {
	v := NewValidator()

	data := testStruct{
		Username: "ValidUser",
		Password: "P1!",
	}

	err := v.Struct(data)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Password")
}

func TestValidator_Invalid_Password_TooLong(t *testing.T) {
	v := NewValidator()

	data := testStruct{
		Username: "ValidUser",
		Password: "A1$" + string(make([]byte, 30)),
	}

	err := v.Struct(data)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Password")
}
