package repository

import (
	"fmt"
	"sso/models"
)

func (r *PostgresRepository) CreateGroup(code, comment string) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (code, comment) VALUES ($1, $2) RETURNING id", GroupTable)
	err := r.db.QueryRow(query, code, comment).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *PostgresRepository) GetGroupByID(id int) (models.Group, error) {
	var group models.Group
	query := fmt.Sprintf("SELECT id, code, comment FROM %s WHERE id = $1", GroupTable)
	err := r.db.Get(&group, query, id)
	if err != nil {
		return group, err
	}
	return group, nil
}

func (r *PostgresRepository) GetGroupByCode(code string) (models.Group, error) {
	var group models.Group
	query := fmt.Sprintf("SELECT id, code, comment FROM %s WHERE code = $1", GroupTable)
	err := r.db.Get(&group, query, code)
	if err != nil {
		return group, err
	}
	return group, nil
}

func (r *PostgresRepository) GetAllGroups() ([]models.Group, error) {
	var groups []models.Group
	query := fmt.Sprintf("SELECT id, code, comment FROM %s ORDER BY code", GroupTable)
	err := r.db.Select(&groups, query)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (r *PostgresRepository) UpdateGroup(id int, group models.Group) error {
	query := fmt.Sprintf("UPDATE %s SET code = $1, comment = $2 WHERE id = $3", GroupTable)
	_, err := r.db.Exec(query, group.Code, group.Comment, id)
	return err
}

func (r *PostgresRepository) DeleteGroup(id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", GroupTable)
	_, err := r.db.Exec(query, id)
	return err
}
