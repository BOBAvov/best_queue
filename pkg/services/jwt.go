package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Генерация JWT токена для пользователя
func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUserByTgName(username)
	if err != nil {
		return "", err
	}
	// Проверка пароля
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password+s.salt)); err != nil {
		return "", errors.New("invalid password")
	}
	// Создание токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserId: int(user.Id),
	})
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

	return int(claims.UserId), nil
}
