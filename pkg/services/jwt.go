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
	IsAdmin bool `json:"is_admin"`
}

// Проверка пароля
func (s *AuthService) VerificationPassword(username, password string) (string, error) {
	const op = "VerificationPassword"
	user, err := s.repo.GetUserByTgName(username)
	if err != nil {
		return "", err
	}
	// Проверка пароля
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}
	// Создание токена
	return s.NewToken(models.AuthUser{
		TgNick:   username,
		Password: password,
	})
}

// Проверка и парсинг JWT токена, возвращает user_id
func (s *AuthService) ParseToken(tokenStr string) (int, bool, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return 0, false, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok || !token.Valid {
		return 0, false, errors.New("token is invalid")
	}

	return claims.UserId, claims.IsAdmin, nil
}

func (s *AuthService) NewToken(userAuth models.AuthUser) (string, error) {
	// Заполняем типизированные claims, чтобы совпадало с ParseToken
	userId, err := s.repo.GetUserIdByTgNick(userAuth.TgNick)
	if err != nil {
		return "", err
	}
	isAdmin, err := s.repo.GetUserIsAdmin(userId)
	if err != nil {
		return "", err
	}
	claims := &tokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
		},
		UserId:  userId,
		IsAdmin: isAdmin,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
