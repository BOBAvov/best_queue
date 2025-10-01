package repository

import (
	"sso/models"

	"github.com/jmoiron/sqlx"
)

const (
	FacultyTable    = "faculties"
	UserTable       = "users"
	GroupTable      = "groups"
	AvailableTable  = "available"
	DepartmentTable = "departments"
	GroupsInQueue   = "groups_in_queue"
	Queues          = "queues"
)

type Repository interface {
	CreateUser(user models.RegisterUser) (int, error)
	GetUserByID(id int) (models.User, error)
	GetUserByTgName(tgName string) (models.User, error)
	CreateGroup(name string) (int, error)
}

func NewRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}
