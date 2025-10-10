// Package repository содержит интерфейсы и структуры для работы с базой данных
package repository

import (
	"sso/models"

	"github.com/jmoiron/sqlx"
)

// Константы с названиями таблиц в базе данных
const (
	UserTable              = "users"              // Таблица пользователей
	GroupTable             = "groups"             // Таблица групп
	QueuesTable            = "queues"             // Таблица очередей
	QueueParticipantsTable = "queue_participants" // Таблица участников очередей
)

// Repository определяет интерфейс для работы с базой данных
type Repository interface {
	// Методы для работы с пользователями
	CreateUser(user models.RegisterUser, groupID int) (int, error) // Создание пользователя
	GetUserByID(id int) (models.User, error)                       // Получение пользователя по ID
	GetUserByTgName(tgName string) (models.User, error)            // Получение пользователя по Telegram никнейму
	GetUserIsAdmin(id int) (bool, error)                           // Проверка статуса администратора
	GetUserIdByTgNick(tgNick string) (int, error)                  // Получение ID пользователя по Telegram никнейму
	GetAllUsers() ([]models.User, error)                           // Получение всех пользователей
	UpdateUser(id int, user models.User) error                     // Обновление пользователя
	DeleteUser(id int) error                                       // Удаление пользователя

	// Методы для работы с группами
	CreateGroup(code, comment string) (int, error)    // Создание группы
	GetGroupByID(id int) (models.Group, error)        // Получение группы по ID
	GetGroupByCode(code string) (models.Group, error) // Получение группы по коду
	GetAllGroups() ([]models.Group, error)            // Получение всех групп
	UpdateGroup(id int, group models.Group) error     // Обновление группы
	DeleteGroup(id int) error                         // Удаление группы

	// Методы для работы с очередями
	CreateQueue(queue models.Queue) (int, error) // Создание очереди
	GetQueueByID(id int) (models.Queue, error)   // Получение очереди по ID
	GetAllQueues() ([]models.Queue, error)       // Получение всех очередей
	UpdateQueue(queue models.Queue) error        // Обновление очереди
	DeleteQueue(id int) error                    // Удаление очереди

	// Методы для работы с участниками очередей
	JoinQueue(queueID, userID int) (int, error)                          // Присоединение к очереди
	LeaveQueue(queueID, userID int) error                                // Покидание очереди
	GetQueueParticipants(queueID int) ([]models.QueueParticipant, error) // Получение участников очереди
	GetUserQueuePosition(queueID, userID int) (int, error)               // Получение позиции пользователя в очереди
	ShiftQueue(queueID int) error                                        // Сдвиг очереди
	GetNextQueuePosition(queueID int) (int, error)                       // Получение следующей позиции в очереди
}

// NewRepository создает новый экземпляр PostgreSQL репозитория
func NewRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}
