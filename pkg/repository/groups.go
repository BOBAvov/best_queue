package repository

import "fmt"

func (r *PostgresRepository) CreateGroup(name string) (id int, err error) {
	query := fmt.Sprintf("INSERT INTO %s (group_name) VALUES ($1)", GroupTable)
	err = r.db.QueryRow(query, name).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *PostgresRepository) GetGroupIdByCode(code string) (id int, err error) {
	query := fmt.Sprintf("SELECT id FROM %s WHERE code=$1", GroupTable)
	err = r.db.QueryRow(query, code).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
