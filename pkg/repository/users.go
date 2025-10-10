package repository

import (
	"fmt"
	"sso/models"
)

func (r *PostgresRepository) CreateUser(user models.RegisterUser, idGroup int) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (username, tg_nick, group_id,password_hash) VALUES ($1, $2, $3,$4 ) RETURNING id", UserTable)
	err := r.db.QueryRow(query, user.Username, user.TgNick, idGroup, user.Password).Scan(&id)
	return id, err
}

// GetAllUsers возвращает всех пользователей
func (r *PostgresRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	query := fmt.Sprintf("SELECT id, username, tg_nick, group_id, password_hash, is_admin FROM %s ORDER BY username", UserTable)
	err := r.db.Select(&users, query)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// UpdateUser обновляет данные пользователя
func (r *PostgresRepository) UpdateUser(id int, user models.User) error {
	query := fmt.Sprintf("UPDATE %s SET username = $1, tg_nick = $2, group_id = $3, is_admin = $4 WHERE id = $5", UserTable)
	_, err := r.db.Exec(query, user.Username, user.TgNick, user.GroupID, user.IsAdmin, id)
	return err
}

// DeleteUser удаляет пользователя
func (r *PostgresRepository) DeleteUser(id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", UserTable)
	_, err := r.db.Exec(query, id)
	return err
}

func (r *PostgresRepository) GetUserByID(id int) (models.User, error) {
	var user models.User
	query := "SELECT id, username, tg_nick, group_id, password_hash, is_admin FROM users WHERE id = $1"
	err := r.db.Get(&user, query, id)
	if err != nil {
		return user, err
	}
	return user, nil
}

// GetUserByTgName возвращает пользователя по Telegram имени
func (r *PostgresRepository) GetUserByTgName(tgName string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT id, username, tg_nick, group_id, password_hash, is_admin FROM %s WHERE tg_nick = $1", UserTable)
	err := r.db.Get(&user, query, tgName)
	if err != nil {
		return user, err
	}
	return user, nil
}

// GetUserIsAdmin проверяет, является ли пользователь администратором
func (r *PostgresRepository) GetUserIsAdmin(id int) (bool, error) {
	var isAdmin bool
	query := fmt.Sprintf("SELECT is_admin FROM %s WHERE id = $1", UserTable)
	err := r.db.Get(&isAdmin, query, id)
	if err != nil {
		return false, err
	}
	return isAdmin, nil
}

// GetUserIdByTgNick возвращает ID пользователя по Telegram нику
func (r *PostgresRepository) GetUserIdByTgNick(tgNick string) (int, error) {
	var id int
	query := fmt.Sprintf("SELECT id FROM %s WHERE tg_nick = $1", UserTable)
	err := r.db.Get(&id, query, tgNick)
	if err != nil {
		return 0, err
	}
	return id, nil
}
