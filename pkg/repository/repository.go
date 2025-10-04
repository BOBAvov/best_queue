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
	// User
	CreateUser(user models.RegisterUser, idGroup int) (int, error)
	GetUserByID(id int) (models.User, error)
	GetUserByTgName(tgName string) (models.User, error)
	GetUserIsAdmin(id int) (bool, error)
	GetUserIdByTgNick(tgNick string) (id int, err error)
	// Group
	CreateGroup(name string) (int, error)
	GetGroupIdByCode(code string) (id int, err error)
	// Faculcy
	CreateFaculty(faculty models.Faculty) (int, error)
	GetFacultyByCode(code string) (models.Faculty, error)
}

func NewRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}
