package services

import (
	"crypto/rand"
	"sso/pkg/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	tokenTTL  = 1000 * time.Hour
	jwtSecret = "my_super_secret_key"
)

// Структура для хранения JWT claims
type tokenClaims struct {
	jwt.RegisteredClaims
	UserId int `json:"user_id"`
}

// Интерфейс авторизации
type Authorization interface {
	CreateUser(user repository.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(tokenStr string) (int, error)
}

// Сервис авторизации
type AuthService struct {
	salt string
	repo repository.Repository
}

// Конструктор сервиса авторизации
func NewAuthService(repo repository.Repository) *AuthService {
	return &AuthService{
		repo: repo,
		salt: generateSalt(),
	}
}

// Создание пользователя с хешированием пароля
func (s *AuthService) CreateUser(user repository.User) (int, error) {
	id_group, err := s.repo.GetIdGroupByName(user.Group)
	if err != nil {
		return 0, err
	}
	hashedPassword, err := s.generatePasswordHash(user.Password)
	if err != nil {
		return 0, err
	}
	user.Password = hashedPassword
	user.GroupID = id_group
	return s.repo.CreateUser(user)
}

func generateSalt() string {
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		panic(err)
	}
	return string(salt)
}

// Хеширование пароля с солью
func (s *AuthService) generatePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password+s.salt), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *AuthService) GetIdGroupByName(groupName string) (int, error) {
	return s.repo.GetIdGroupByName(groupName)
}
