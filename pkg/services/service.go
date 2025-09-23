package services

import (
	"errors"
	"sso/pkg/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	salt      = "crwty3146738yegy68td2f76287382893293201328321br63b26886xf7z36z"
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
	repo repository.Authorization
}

// Конструктор сервиса авторизации
func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

// Создание пользователя с хешированием пароля
func (s *AuthService) CreateUser(user repository.User) (int, error) {
	hashedPassword, err := generatePasswordHash(user.Password)
	if err != nil {
		return 0, err
	}
	user.Password = hashedPassword
	return s.repo.CreateUser(user)
}

// Генерация JWT токена для пользователя
func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	// Проверка пароля
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password+salt)); err != nil {
		return "", err
	}

	// Создание токена с user_id
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserId: user.Id,
	})

	// Возвращаем строку токена
	return token.SignedString([]byte(jwtSecret))
}

// Проверка и парсинг JWT токена, возвращает user_id
func (s *AuthService) ParseToken(tokenStr string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok || !token.Valid {
		return 0, errors.New("token is invalid")
	}

	return claims.UserId, nil
}

// Хеширование пароля с солью
func generatePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
	return string(bytes), err
}
