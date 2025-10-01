package services

import (
	"sso/models"
	"sso/pkg/repository"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	tokenTTL  = 1000 * time.Hour
	jwtSecret = "my_super_secret_key"
)

// Структура для хранения JWT claims

// Интерфейс авторизации
type Authorization interface {
	CreateUser(user models.RegisterUser) (int, error)
	//GenerateToken(username, password string) (string, error)
	ParseToken(tokenStr string) (int, error)
	CreateGroup(user string) (int, error)
}

// Сервис авторизации
type AuthService struct {
	repo repository.Repository
}

// Конструктор сервиса авторизации
func NewAuthService(repo repository.Repository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

// Создание пользователя с хешированием пароля
func (s *AuthService) CreateUser(user models.RegisterUser) (int, error) {
	hashedPassword, err := s.generatePasswordHash(user.Password)
	if err != nil {
		return 0, err
	}
	user.Password = hashedPassword
	return s.repo.CreateUser(user)
}

//

// Хеширование пароля с солью
func (s *AuthService) generatePasswordHash(password string) (string, error) {
	// Используем соль, встроенную в bcrypt, без дополнительной пользовательской соли
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *AuthService) CreateGroup(name string) (int, error) {
	return s.repo.CreateGroup(name)
}
