package repository

import (
	"fmt"
	"sso/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
