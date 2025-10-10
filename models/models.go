package models

import (
	"database/sql"
	"time"
)

type Config struct {
	Port string
	DB   DBConfig
}

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type Faculty struct {
	ID       int            `db:"id" json:"id"`
	Code     string         `db:"code" json:"code" binding:"required"`
	Name     string         `db:"name" json:"name" binding:"required"`
	Comments sql.NullString `db:"comments" json:"comments"`
}

// Department соответствует таблице "Departments"
type Department struct {
	ID        int            `db:"id" json:"id"`
	Code      string         `db:"code" json:"code"`
	Name      string         `db:"name" json:"name"`
	FacultyID int            `db:"faculty_id" json:"faculty_id"`
	Comment   sql.NullString `db:"comment" json:"comment"`
}

// Group соответствует таблице "Groups"
type Group struct {
	ID           int            `db:"id" json:"id"`
	Code         string         `db:"code" json:"code" binding:"required"`
	Name         string         `db:"name" json:"name" binding:"required"`
	DepartmentID int            `db:"department_id" json:"department_id"`
	Comment      sql.NullString `db:"comment" json:"comment"`
}

type RegisterUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	TgNick   string `json:"tg_nick"`
	Group    string `json:"group"`
}

// User соответствует таблице "Users"
type User struct {
	ID           int    `db:"id" json:"id"`
	Username     string `db:"username" json:"username"`
	TgNick       string `db:"tg_nick" json:"tg_nick"`
	GroupID      int    `db:"group_id" json:"group_id"`
	PasswordHash string `db:"password_hash" json:"-"`
	IsAdmin      bool   `db:"is_admin" json:"is_admin"`
}

type Available struct {
	ID   int    `db:"id" json:"id"`
	Type string `db:"type" json:"type"`
}

type Queue struct {
	ID          int       `db:"id" json:"id"`
	Tittle      string    `db:"title" json:"title" binding:"required"`
	GroupName   string    `db:"-" json:"group_name" binding:"required"`
	GroupID     int       `db:"group_id" json:"group_id"`
	AvailableID int       `db:"available_id" json:"available_id"`
	TimeStart   time.Time `db:"time_start" json:"time_start"`
	TimeAdd     int       `json:"time_add"`
	TimeEnd     time.Time `db:"time_end" json:"time_end"`
}

type GroupInQueue struct {
	ID      int64 `db:"id" json:"id"`
	QueueID int64 `db:"queue_id" json:"queue_id"`
	GroupID int64 `db:"group_id" json:"group_id"`
}

type AuthUser struct {
	TgNick   string `json:"tg_nick" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type QueueParticipant struct {
	ID       int       `db:"id" json:"id"`
	QueueID  int       `db:"queue_id" json:"queue_id"`
	UserID   int       `db:"user_id" json:"user_id"`
	Position int       `db:"position" json:"position"`
	JoinedAt time.Time `db:"joined_at" json:"joined_at"`
	IsActive bool      `db:"is_active" json:"is_active"`
}

type CreateQueueRequest struct {
	Title     string `json:"title" binding:"required"`
	GroupCode string `json:"group_code" binding:"required"`
	TimeAdd   int    `json:"time_add"`
}

type JoinQueueRequest struct {
	QueueID int `json:"queue_id" binding:"required"`
}

type FacultyCreateRequest struct {
	Code     string `json:"code" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Comments string `json:"comments"`
}
