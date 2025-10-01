package main

import (
	"bufio"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
	"sso/models"
	"sso/pkg/config"
	"sso/pkg/repository"
	"strings"
)

type Data struct {
	db *sqlx.DB
}

func main() {
	cfg := config.LoadConfig()
	db, err := repository.NewPostgresDB(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("DB closed")
	}()

	faculty, err := ParseTxt("facultacy.txt")
	if err != nil {
		log.Fatal(err)
	}
	data := NewData(db)
	for _, t := range faculty {
		id, err := data.CreateFaculty(models.Faculty{
			Name: t[1],
			Code: t[0],
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id)

	}

}
func NewData(db *sqlx.DB) *Data {
	return &Data{db: db}
}

func (h *Data) AutoMigrate() {

}

func (h *Data) CreateFaculty(faculty models.Faculty) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (code,name,comments) VALUES ($1,$2,$3) RETURNING id", "faculties")
	err := h.db.QueryRow(query, faculty.Code, faculty.Name, faculty.Comments).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (h *Data) GetFacultyByCode(code string) (models.Faculty, error) {
	var faculty models.Faculty

}

func ParseTxt(fileName string) ([][2]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var faculties [][2]string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Разделяем строку по первому пробелу
		parts := strings.SplitN(line, " ", 2)
		if len(parts) == 2 {
			faculties = append(faculties, [2]string{parts[0], parts[1]})
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return faculties, nil
}
