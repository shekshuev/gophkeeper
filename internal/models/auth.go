package models

// LoginUserDTO представляет данные, передаваемые пользователем при попытке входа в систему.
type LoginUserDTO struct {
	UserName string `json:"user_name" validate:"required,min=5,max=30,alphanumunderscore,startswithalpha"` // Логин: от 5 до 30 символов, буквы/цифры/подчёркивание, начинается с буквы
	Password string `json:"password" validate:"required,password"`                                         // Пароль: обязательный, соответствует пользовательским правилам
}

// RegisterUserDTO используется при регистрации нового пользователя.
// Включает подтверждение пароля и валидацию имени/фамилии.
type RegisterUserDTO struct {
	UserName        string `json:"user_name" validate:"required,min=5,max=30,alphanumunderscore,startswithalpha"` // Логин
	Password        string `json:"password" validate:"required,password"`                                         // Пароль
	PasswordConfirm string `json:"password_confirm" validate:"required,password,eqfield=Password"`                // Подтверждение пароля (должно совпадать с Password)
	FirstName       string `json:"first_name" validate:"required,min=1,max=30,alphaunicode"`                      // Имя: только буквы (включая Unicode)
	LastName        string `json:"last_name" validate:"required,min=1,max=30,alphaunicode"`                       // Фамилия: только буквы (включая Unicode)
}

// ReadTokenDTO содержит access и refresh токены, возвращаемые после успешной аутентификации.
type ReadTokenDTO struct {
	AccessToken  string `json:"access_token"`  // JWT access token (короткоживущий)
	RefreshToken string `json:"refresh_token"` // JWT refresh token (для обновления access токена)
}
