package repository

import (
	"fmt"
	"sso/models"
)

func (r *PostgresRepository) GetUserByTgName(tgName string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT id, tg_name, password_hash, is_admin, group_name FROM %s WHERE tg_name=$1", UserTable)
	err := r.db.Get(&user, query, tgName)
	return user, err
}

func (r *PostgresRepository) GetUserByID(id int) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT id, tg_name, password_hash, is_admin, group_name FROM %s WHERE id=$1", UserTable)
	err := r.db.Get(&user, query, id)
	return user, err
}

func (r *PostgresRepository) CreateUser(user models.RegisterUser) (int, error) {
	// Добавление пользователя возможно только после проверки Существования Группы

	var id int
	query := fmt.Sprintf("INSERT INTO %s (username, tg_nick, group_id, password_hash) VALUES ($1,$2,$3,$4) RETURNING id", UserTable)
	err := r.db.QueryRow(query, user.Username, user.TgNick, user.Group, user.Password).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
