// Package services содержит бизнес-логику приложения
package services

import (
	"sso/models"
	"sso/pkg/repository"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Константы для JWT токенов
const (
	tokenTTL  = 1000 * time.Hour      // Время жизни JWT токена
	jwtSecret = "my_super_secret_key" // Секретный ключ для подписи JWT токенов
)

// Authorization определяет интерфейс для работы с авторизацией и управлением данными
type Authorization interface {
	// Аутентификация и авторизация
	CreateUser(user models.RegisterUser) (int, error) // Создание нового пользователя
	NewToken(user models.AuthUser) (string, error)    // Генерация JWT токена
	ParseToken(tokenStr string) (int, bool, error)    // Парсинг JWT токена

	// Управление группами
	CreateGroup(code, comment string) (int, error)    // Создание новой группы
	GetGroupByID(id int) (models.Group, error)        // Получение группы по ID
	GetGroupByCode(code string) (models.Group, error) // Получение группы по коду
	GetAllGroups() ([]models.Group, error)            // Получение всех групп
	UpdateGroup(id int, group models.Group) error     // Обновление группы
	DeleteGroup(id int) error                         // Удаление группы

	// Управление очередями
	CreateQueue(queue models.Queue) (int, error) // Создание новой очереди
	GetQueueByID(id int) (models.Queue, error)   // Получение очереди по ID
	GetAllQueues() ([]models.Queue, error)       // Получение всех очередей
	UpdateQueue(queue models.Queue) error        // Обновление очереди
	DeleteQueue(id int) error                    // Удаление очереди

	// Управление участниками очередей
	JoinQueue(queueID, userID int) (int, error)                          // Присоединение к очереди
	LeaveQueue(queueID, userID int) error                                // Покидание очереди
	GetQueueParticipants(queueID int) ([]models.QueueParticipant, error) // Получение участников очереди
	ShiftQueue(queueID int) error                                        // Сдвиг очереди

	// Управление пользователями
	GetUserByID(id int) (models.User, error)   // Получение пользователя по ID
	GetAllUsers() ([]models.User, error)       // Получение всех пользователей
	UpdateUser(id int, user models.User) error // Обновление пользователя
	DeleteUser(id int) error                   // Удаление пользователя
}

// AuthService реализует интерфейс Authorization и содержит бизнес-логику приложения
type AuthService struct {
	repo repository.Repository // Репозиторий для работы с базой данных
}

// NewAuthService создает новый экземпляр сервиса авторизации
func NewAuthService(repo repository.Repository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

// generatePasswordHash хеширует пароль с использованием bcrypt
func (s *AuthService) generatePasswordHash(password string) (string, error) {
	// Используем соль, встроенную в bcrypt, без дополнительной пользовательской соли
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Методы управления группами

// CreateGroup создает новую группу
func (s *AuthService) CreateGroup(code, comment string) (int, error) {
	return s.repo.CreateGroup(code, comment)
}

// GetGroupByID получает группу по ID
func (s *AuthService) GetGroupByID(id int) (models.Group, error) {
	return s.repo.GetGroupByID(id)
}

// GetGroupByCode получает группу по коду
func (s *AuthService) GetGroupByCode(code string) (models.Group, error) {
	return s.repo.GetGroupByCode(code)
}

// GetAllGroups получает все группы
func (s *AuthService) GetAllGroups() ([]models.Group, error) {
	return s.repo.GetAllGroups()
}

// UpdateGroup обновляет информацию о группе
func (s *AuthService) UpdateGroup(id int, group models.Group) error {
	return s.repo.UpdateGroup(id, group)
}

// DeleteGroup удаляет группу
func (s *AuthService) DeleteGroup(id int) error {
	return s.repo.DeleteGroup(id)
}
