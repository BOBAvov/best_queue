// Package services содержит JWT токены и аутентификацию
package services

import (
	"errors"
	"sso/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// tokenClaims представляет структуру JWT токена с пользовательскими данными
type tokenClaims struct {
	jwt.RegisteredClaims
	UserId  int  `json:"user_id"`  // ID пользователя
	IsAdmin bool `json:"is_admin"` // Флаг администратора
}

// VerificationPassword проверяет пароль пользователя и возвращает JWT токен
func (s *AuthService) VerificationPassword(username, password string) (string, error) {
	const op = "VerificationPassword"
	user, err := s.repo.GetUserByTgName(username)
	if err != nil {
		return "", err
	}

	// Проверка пароля с использованием bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}

	// Создание JWT токена при успешной аутентификации
	return s.NewToken(models.AuthUser{
		TgNick:   username,
		Password: password,
	})
}

// ParseToken парсит JWT токен и возвращает ID пользователя и статус администратора
func (s *AuthService) ParseToken(tokenStr string) (int, bool, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return 0, false, err
	}

	// Извлекаем claims из токена
	claims, ok := token.Claims.(*tokenClaims)
	if !ok || !token.Valid {
		return 0, false, errors.New("token is invalid")
	}

	return claims.UserId, claims.IsAdmin, nil
}

// NewToken создает новый JWT токен для пользователя
func (s *AuthService) NewToken(userAuth models.AuthUser) (string, error) {
	// Получаем ID пользователя по Telegram никнейму
	userId, err := s.repo.GetUserIdByTgNick(userAuth.TgNick)
	if err != nil {
		return "", err
	}

	// Получаем статус администратора
	isAdmin, err := s.repo.GetUserIsAdmin(userId)
	if err != nil {
		return "", err
	}

	// Создаем claims для JWT токена
	claims := &tokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)), // Устанавливаем время истечения
		},
		UserId:  userId,
		IsAdmin: isAdmin,
	}

	// Создаем и подписываем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
