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

func (h *PostgresRepository) GetFacultyByCode(code string) (models.Faculty, error) {
	var faculty models.Faculty
	query := fmt.Sprintf("SELECT id FROM %s WHERE code=$1", FacultyTable)
	err := h.db.Get(&faculty, query, code)
	if err != nil {
		return faculty, err
	}
	return faculty, nil
}
