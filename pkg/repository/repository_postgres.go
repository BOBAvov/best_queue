// Package repository содержит реализацию PostgreSQL репозитория
package repository

import (
	"fmt"
	"sso/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// PostgresRepository реализует интерфейс Repository для PostgreSQL
type PostgresRepository struct {
	db *sqlx.DB // Подключение к базе данных PostgreSQL
}

// NewPostgresDB создает новое подключение к PostgreSQL базе данных
func NewPostgresDB(cfg models.DBConfig) (*sqlx.DB, error) {
	// Формируем строку подключения к PostgreSQL
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.Username, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Проверяем подключение к базе данных
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Close закрывает подключение к базе данных
func (r *PostgresRepository) Close() error {
	return r.db.Close()
}
