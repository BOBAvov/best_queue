package repository

import "fmt"

func (r *PostgresRepository) CreateGroup(code, name string) (id int, err error) {
	query := fmt.Sprintf("INSERT INTO %s (code, name) VALUES ($1, $2) RETURNING id", GroupTable)
	err = r.db.QueryRow(query, code, name).Scan(&id)
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
