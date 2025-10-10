// Package models содержит структуры данных для HTTP SSO API
package models

import (
	"database/sql"
	"time"
)

// Config представляет конфигурацию приложения
type Config struct {
	Port string   // Порт для запуска HTTP сервера
	DB   DBConfig // Конфигурация базы данных
}

// DBConfig содержит параметры подключения к базе данных PostgreSQL
type DBConfig struct {
	Host     string // Хост базы данных
	Port     string // Порт базы данных
	Username string // Имя пользователя БД
	Password string // Пароль пользователя БД
	DBName   string // Имя базы данных
	SSLMode  string // Режим SSL подключения
}

// CreateGroupRequest представляет запрос на создание новой группы
type CreateGroupRequest struct {
	Code    string `json:"code" binding:"required"` // Код группы (обязательное поле)
	Comment string `json:"comment"`                 // Комментарий к группе
}

// Group представляет группу студентов, соответствует таблице "Groups" в БД
type Group struct {
	ID      int            `db:"id" json:"id"`                        // Уникальный идентификатор группы
	Code    string         `db:"code" json:"code" binding:"required"` // Код группы (например, "ИУ7-12Б")
	Comment sql.NullString `db:"comment" json:"comment"`              // Описание группы (может быть NULL)
}

// RegisterUser представляет данные для регистрации нового пользователя
type RegisterUser struct {
	Username string `json:"username"` // Имя пользователя
	Password string `json:"password"` // Пароль пользователя
	TgNick   string `json:"tg_nick"`  // Telegram никнейм
	Group    string `json:"group"`    // Код группы пользователя
}

// User представляет пользователя системы, соответствует таблице "Users" в БД
type User struct {
	ID           int    `db:"id" json:"id"`             // Уникальный идентификатор пользователя
	Username     string `db:"username" json:"username"` // Имя пользователя
	TgNick       string `db:"tg_nick" json:"tg_nick"`   // Telegram никнейм
	GroupID      int    `db:"group_id" json:"group_id"` // ID группы пользователя
	PasswordHash string `db:"password_hash" json:"-"`   // Хеш пароля (не возвращается в JSON)
	IsAdmin      bool   `db:"is_admin" json:"is_admin"` // Флаг администратора
}

// Queue представляет очередь на консультацию, соответствует таблице "Queues" в БД
type Queue struct {
	ID        int       `db:"id" json:"id"`                          // Уникальный идентификатор очереди
	Title     string    `db:"title" json:"title" binding:"required"` // Название очереди
	TimeStart time.Time `db:"time_start" json:"time_start"`          // Время начала приема
	TimeEnd   time.Time `db:"time_end" json:"time_end"`              // Время окончания приема
}

// AuthUser представляет данные для аутентификации пользователя
type AuthUser struct {
	TgNick   string `json:"tg_nick" binding:"required"`  // Telegram никнейм (обязательное поле)
	Password string `json:"password" binding:"required"` // Пароль пользователя (обязательное поле)
}

// QueueParticipant представляет участника очереди, соответствует таблице "QueueParticipants" в БД
type QueueParticipant struct {
	ID       int       `db:"id" json:"id"`               // Уникальный идентификатор участника
	QueueID  int       `db:"queue_id" json:"queue_id"`   // ID очереди
	UserID   int       `db:"user_id" json:"user_id"`     // ID пользователя
	Position int       `db:"position" json:"position"`   // Позиция в очереди
	JoinedAt time.Time `db:"joined_at" json:"joined_at"` // Время присоединения к очереди
	IsActive bool      `db:"is_active" json:"is_active"` // Флаг активности участника
}

// CreateQueueRequest представляет запрос на создание новой очереди
type CreateQueueRequest struct {
	Title     string    `json:"title" binding:"required"`      // Название очереди (обязательное поле)
	TimeStart time.Time `json:"time_start" binding:"required"` // Время начала приема (обязательное поле)
	TimeEnd   time.Time `json:"time_end" binding:"required"`   // Время окончания приема (обязательное поле)
}

// JoinQueueRequest представляет запрос на присоединение к очереди
type JoinQueueRequest struct {
	QueueID int `json:"queue_id" binding:"required"` // ID очереди для присоединения (обязательное поле)
}

// GroupCreateRequest представляет запрос на создание новой группы (дублирует CreateGroupRequest)
type GroupCreateRequest struct {
	Code    string `json:"code" binding:"required"` // Код группы (обязательное поле)
	Comment string `json:"comment"`                 // Комментарий к группе
}
