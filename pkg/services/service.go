package services

import (
	"errors"
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
	NewToken(user models.AuthUser) (string, error)
	ParseToken(tokenStr string) (int, bool, error)
	CreateGroup(user string) (int, error)
	CreateQueue(queue models.Queue) (int, error)
	GetQueueById(id string) (models.Queue, error)
	GetAllQueues() ([]models.Queue, error)
	UpdateQueue(queue models.Queue) error
	DeleteQueue(id string) error
	// Queue Participants
	JoinQueue(queueID, userID int) (int, error)
	LeaveQueue(queueID, userID int) error
	GetQueueParticipants(queueID int) ([]models.QueueParticipant, error)
	ShiftQueue(queueID int) error
	// User management
	GetUserByID(id int) (models.User, error)
	GetAllUsers() ([]models.User, error)
	UpdateUser(id int, user models.User) error
	DeleteUser(id int) error
	GetGroupIdByCode(code string) (int, error)
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

// Хеширование пароля с солью
func (s *AuthService) generatePasswordHash(password string) (string, error) {
	// Используем соль, встроенную в bcrypt, без дополнительной пользовательской соли
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *AuthService) CreateGroup(name string) (int, error) {
	code := ""
	for i, ch := range name {
		if ch == '-' {
			code = name[0:i]
			break
		}
	}
	if code == "" {
		return 0, errors.New("Group Name is required")
	}

	return s.repo.CreateGroup(name, code)
}
