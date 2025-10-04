package repository

import (
	"fmt"
	"log"
	"sso/models"
)

func (r *PostgresRepository) GetUserByTgName(tgName string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE tg_name=$1", UserTable)
	err := r.db.Get(&user, query, tgName)
	return user, err
}

func (r *PostgresRepository) GetUserByID(id int) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", UserTable)
	err := r.db.Get(&user, query, id)
	return user, err
}

func (r *PostgresRepository) CreateUser(user models.RegisterUser, idGroup int) (idUser int, err error) {
	query := fmt.Sprintf("INSERT INTO %s (username, tg_nick, group_id, password_hash) VALUES ($1,$2,$3,$4) RETURNING id", UserTable)
	err = r.db.QueryRow(query, user.Username, user.TgNick, idGroup, user.Password).Scan(&idUser)
	if err != nil {
		return 0, err
	}
	return idUser, nil
}

func (r *PostgresRepository) GetUserIdByTgNick(tgNick string) (id int, err error) {
	const op = "GetUserIdByTgNick"
	query := fmt.Sprintf("SELECT id FROM %s WHERE tg_nick=$1", UserTable)
	err = r.db.Get(&id, query, tgNick)
	if err != nil {
		log.Printf("%s: %v", op, err)
		return 0, err
	}
	return id, err
}

func (r *PostgresRepository) GetUserIsAdmin(id int) (bool, error) {
	var isAdmin bool
	query := fmt.Sprintf("SELECT is_admin FROM %s WHERE id=$1", UserTable)
	err := r.db.Get(&isAdmin, query, id)
	return isAdmin, err
}
