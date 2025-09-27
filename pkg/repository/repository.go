package repository

import (
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	CreateUser(user User) (int, error)
	GetUserByID(id int) (User, error)
	GetUserByTgName(tgName string) (User, error)
	GetIdGroupByName(groupName string) (int, error)
}

func NewRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}
