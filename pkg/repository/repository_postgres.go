package repository

import (
	"fmt"
	"sso/models"

	"github.com/jmoiron/sqlx"
)

const (
	usersTable  = "Users"
	groupsTable = "Groups"
)

type PostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresDB(cfg models.DBConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.Username, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (r *PostgresRepository) Close() error {
	return r.db.Close()
}
func (r *PostgresRepository) CreateUser(user User) (int, error) {
	// Добавление пользователя возможно только после проверки Существования Группы

	var id int
	query := fmt.Sprintf("INSERT INTO %s (username, tg_nick, group_id, password_hash) VALUES ($1,$2,$3,$4) RETURNING id", usersTable)
	err := r.db.QueryRow(query, user.Username, user.Tg_nick, user.GroupID, user.Password).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *PostgresRepository) GetUserByID(id int) (User, error) {
	var user User
	query := fmt.Sprintf("SELECT id, tg_name, password_hash, is_admin, group_name FROM %s WHERE id=$1", usersTable)
	err := r.db.Get(&user, query, id)
	return user, err
}
func (r *PostgresRepository) GetUserByTgName(tgName string) (User, error) {
	var user User
	query := fmt.Sprintf("SELECT id, tg_name, password_hash, is_admin, group_name FROM %s WHERE tg_name=$1", usersTable)
	err := r.db.Get(&user, query, tgName)
	return user, err
}

func (r *PostgresRepository) GetIdGroupByName(groupName string) (int, error) {
	var id int
	query := fmt.Sprintf("SELECT id FROM %s WHERE group_name=$1", groupsTable)
	err := r.db.Get(&id, query, groupName)

	if err != nil {
		return 0, err
	}

	return id, nil
}
