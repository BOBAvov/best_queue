package services

import (
	"crypto/rand"
	"errors"
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
	repo repository.Authorization
}

// Конструктор сервиса авторизации
func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{
		repo: repo,
		salt: generateSalt(),
	}
}

// Создание пользователя с хешированием пароля
func (s *AuthService) CreateUser(user repository.User) (int, error) {
	hashedPassword, err := s.generatePasswordHash(user.Password)
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
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password+s.salt)); err != nil {
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
