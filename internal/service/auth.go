package service

import (
	"context"
	"strconv"

	"github.com/shekshuev/gophkeeper/internal/config"
	"github.com/shekshuev/gophkeeper/internal/models"
	"github.com/shekshuev/gophkeeper/internal/repository"
	"github.com/shekshuev/gophkeeper/internal/utils"
)

// AuthServiceImpl — реализация интерфейса AuthService.
// Отвечает за логику регистрации, аутентификации и генерации JWT-токенов.
type AuthServiceImpl struct {
	repo repository.UserRepository // Репозиторий пользователей
	cfg  *config.Config            // Конфигурация приложения (секреты и срок жизни токенов)
}

// NewAuthServiceImpl создаёт новый экземпляр AuthServiceImpl с указанным репозиторием и конфигурацией.
func NewAuthServiceImpl(repo repository.UserRepository, cfg *config.Config) *AuthServiceImpl {
	return &AuthServiceImpl{repo: repo, cfg: cfg}
}

// Login выполняет аутентификацию пользователя по логину и паролю.
// При успехе возвращает пару access/refresh токенов.
func (s *AuthServiceImpl) Login(ctx context.Context, dto models.LoginUserDTO) (*models.ReadTokenDTO, error) {
	user, err := s.repo.GetUserByUserName(ctx, dto.UserName)
	if err != nil {
		return nil, ErrUserNotFound
	}
	if !utils.VerifyPassword(dto.Password, user.PasswordHash) {
		return nil, ErrWrongPassword
	}
	return s.generateTokenPair(*user)
}

// Register регистрирует нового пользователя и сразу возвращает access/refresh токены.
// Пароль хешируется перед сохранением.
func (s *AuthServiceImpl) Register(ctx context.Context, dto models.RegisterUserDTO) (*models.ReadTokenDTO, error) {
	createDTO := models.CreateUserDTO{
		UserName:     dto.UserName,
		PasswordHash: utils.HashPassword(dto.Password),
		FirstName:    dto.FirstName,
		LastName:     dto.LastName,
	}
	user, err := s.repo.CreateUser(ctx, createDTO)
	if err != nil {
		return nil, err
	}
	return s.generateTokenPair(*user)
}

// generateTokenPair создаёт access и refresh JWT-токены для пользователя.
// Токены подписываются соответствующими секретами из конфигурации.
func (s *AuthServiceImpl) generateTokenPair(user models.ReadAuthUserDataDTO) (*models.ReadTokenDTO, error) {
	userID := strconv.FormatUint(user.ID, 10)

	accessToken, err := utils.CreateToken(
		s.cfg.AccessTokenSecret,
		userID,
		s.cfg.AccessTokenExpires,
	)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.CreateToken(
		s.cfg.RefreshTokenSecret,
		userID,
		s.cfg.RefreshTokenExpires,
	)
	if err != nil {
		return nil, err
	}

	return &models.ReadTokenDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
