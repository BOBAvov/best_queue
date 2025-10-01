package services

import (
	"errors"
	"sso/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type tokenClaims struct {
	jwt.RegisteredClaims
	UserId  int  `json:"user_id"`
	Isadmin bool `json:"is_admin"`
}

// Проверка пароля
func (s *AuthService) VerificationPassword(username, password string) (string, error) {
	user, err := s.repo.GetUserByTgName(username)
	if err != nil {
		return "", err
	}
	// Проверка пароля
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}
	// Создание токена
	return NewToken(user)
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

func NewToken(user models.User) (string, error) {
	// Заполняем типизированные claims, чтобы совпадало с ParseToken
	claims := &tokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
		},
		UserId:  int(user.ID),
		Isadmin: user.IsAdmin,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
