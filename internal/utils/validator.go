package utils

import (
	"log"

	"github.com/dlclark/regexp2"
	"github.com/go-playground/validator/v10"
)

// NewValidator инициализирует и возвращает валидатор с зарегистрированными кастомными правилами:
// - alphanumunderscore: допускаются только буквы, цифры и подчёркивания
// - startswithalpha: значение должно начинаться с буквы
// - password: сложный пароль (буквы + цифры + спецсимволы, от 5 до 30 символов)
func NewValidator() *validator.Validate {
	validate := validator.New()

	err := validate.RegisterValidation("alphanumunderscore", alphanumUnderscore)
	if err != nil {
		log.Fatalf("Error registering validator 'alphanumunderscore': %v", err)
	}

	err = validate.RegisterValidation("startswithalpha", startsWithAlpha)
	if err != nil {
		log.Fatalf("Error registering validator 'startswithalpha': %v", err)
	}

	err = validate.RegisterValidation("password", passwordValidation)
	if err != nil {
		log.Fatalf("Error registering validator 'password': %v", err)
	}

	return validate
}

// alphanumUnderscore проверяет, что значение состоит только из букв, цифр и подчёркиваний.
func alphanumUnderscore(fl validator.FieldLevel) bool {
	return regexValidation(fl.Field().String(), `^[a-zA-Z0-9_]+$`)
}

// startsWithAlpha проверяет, что значение начинается с буквы (не цифрой).
func startsWithAlpha(fl validator.FieldLevel) bool {
	return regexValidation(fl.Field().String(), `^[A-Za-z]`)
}

// passwordValidation проверяет, что пароль:
// - от 5 до 30 символов,
// - содержит хотя бы одну букву,
// - хотя бы одну цифру,
// - хотя бы один спецсимвол из набора @$!%*?&
func passwordValidation(fl validator.FieldLevel) bool {
	return regexValidation(fl.Field().String(), `^(?=.*[a-zA-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{5,30}$`)
}

// regexValidation выполняет проверку строки по регулярному выражению, используя regexp2.
func regexValidation(field, regex string) bool {
	re := regexp2.MustCompile(regex, regexp2.None)
	matched, _ := re.MatchString(field)
	return matched
}
