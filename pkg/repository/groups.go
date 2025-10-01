package repository

import "fmt"

func (r *PostgresRepository) CreateGroup(name string) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (group_name) VALUES ($1)", GroupTable)
	err := r.db.QueryRow(query, name).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
