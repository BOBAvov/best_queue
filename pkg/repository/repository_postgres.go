package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	usersTable = "users"
)

type PostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port))
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

func (r *PostgresRepository) CreateUser(user User) (int32, error) {
	var id int32
	query := fmt.Sprintf("INSERT INTO %s (tg_name,password_hash,is_admin,group_name) VALUES ($1,$2,$3,$4) RETURNING id", usersTable)
	row := r.db.QueryRow(query, user.Tg_name, user.Password, user.IsAdmin, user.Group)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *PostgresRepository) GetUserByID(id int32) (User, error) {
	var user User
	query := fmt.Sprintf("SELECT tg_name FROM %s WHERE id=$1", usersTable)
	row := r.db.QueryRow(query, id)
	if err := row.Scan(&user.Tg_name); err != nil {
		return User{}, err
	}
	return user, nil
}
