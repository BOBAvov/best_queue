package services

import (
	"log"
	"sso/models"
)

// Создание пользователя с хешированием пароля
func (s *AuthService) CreateUser(user models.RegisterUser) (int, error) {
	const op = "CreateUser"
	hashedPassword, err := s.generatePasswordHash(user.Password)
	if err != nil {
		log.Printf("%s: %v", op, err)
		return 0, err
	}
	user.Password = hashedPassword

	var idGroup int
	idGroup, err = s.repo.GetGroupIdByCode(user.Group)
	if err != nil {
		log.Printf("%s: %v", op, err)
		return 0, err
	}

	return s.repo.CreateUser(user, idGroup)
}
