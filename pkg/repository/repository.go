package repository

import (
	"sso/models"

	"github.com/jmoiron/sqlx"
)

const (
	FacultyTable           = "faculties"
	UserTable              = "users"
	GroupTable             = "groups"
	AvailableTable         = "available"
	DepartmentTable        = "departments"
	GroupsInQueue          = "groups_in_queue"
	QueuesTable            = "queues"
	QueueParticipantsTable = "queue_participants"
)

type Repository interface {
	// User
	CreateUser(user models.RegisterUser, idGroup int) (int, error)
	GetUserByID(id int) (models.User, error)
	GetUserByTgName(tgName string) (models.User, error)
	GetUserIsAdmin(id int) (bool, error)
	GetUserIdByTgNick(tgNick string) (id int, err error)
	// Group
	CreateGroup(name string, code string) (int, error)
	GetGroupIdByCode(code string) (id int, err error)
	// Queue
	CreateQueue(queue models.Queue) (int, error)
	GetQueueById(id string) (models.Queue, error)
	GetAllQueues() ([]models.Queue, error)
	UpdateQueue(queue models.Queue) error
	DeleteQueue(id string) error
	// Queue Participants
	JoinQueue(queueID, userID int) (int, error)
	LeaveQueue(queueID, userID int) error
	GetQueueParticipants(queueID int) ([]models.QueueParticipant, error)
	GetUserQueuePosition(queueID, userID int) (int, error)
	ShiftQueue(queueID int) error
	GetNextQueuePosition(queueID int) (int, error)
	// User management
	GetAllUsers() ([]models.User, error)
	UpdateUser(id int, user models.User) error
	DeleteUser(id int) error
}

func NewRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}
