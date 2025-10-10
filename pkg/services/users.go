// Package services содержит методы для работы с пользователями
package services

import (
	"log"
	"sso/models"
)

// CreateUser создает нового пользователя с хешированием пароля
func (s *AuthService) CreateUser(user models.RegisterUser) (int, error) {
	const op = "CreateUser"

	// Хешируем пароль перед сохранением в базу данных
	hashedPassword, err := s.generatePasswordHash(user.Password)
	if err != nil {
		log.Printf("%s: %v", op, err)
		return 0, err
	}
	user.Password = hashedPassword

	// Получаем группу по коду
	group, err := s.repo.GetGroupByCode(user.Group)
	if err != nil {
		log.Printf("%s: %v", op, err)
		return 0, err
	}

	// Создаем пользователя в базе данных
	return s.repo.CreateUser(user, group.ID)
}
