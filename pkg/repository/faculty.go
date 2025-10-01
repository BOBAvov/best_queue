package repository

import (
	"fmt"
	"sso/models"
)

func (h *PostgresRepository) CreateFaculty(faculty models.Faculty) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (code,name,comments) VALUES ($1,$2,$3) RETURNING id", "public.\"Faculties\"")
	err := h.db.QueryRow(query, faculty.Code, faculty.Name, faculty.Comments).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
