package repository

import (
	"fmt"
	"sso/models"
)

func (s *PostgresRepository) CreateQueue(queue models.Queue) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (title, available_id, time_start, time_end) VALUES ($1, $2, $3, $4) RETURNING id", QueuesTable)
	err := s.db.QueryRow(query, queue.Title, queue.AvailableID, queue.TimeStart, queue.TimeEnd).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *PostgresRepository) GetQueueById(id string) (models.Queue, error) {
	var queue models.Queue
	query := fmt.Sprintf("SELECT id, title, available_id, time_start, time_end FROM %s WHERE id = $1", QueuesTable)
	err := s.db.Get(&queue, query, id)
	if err != nil {
		return queue, err
	}
	return queue, nil
}

func (s *PostgresRepository) GetAllQueues() ([]models.Queue, error) {
	var queues []models.Queue
	query := fmt.Sprintf("SELECT id, title, available_id, time_start, time_end FROM %s ORDER BY time_start DESC", QueuesTable)
	err := s.db.Select(&queues, query)
	if err != nil {
		return nil, err
	}
	return queues, nil
}

func (s *PostgresRepository) UpdateQueue(queue models.Queue) error {
	query := fmt.Sprintf("UPDATE %s SET title = $1, available_id = $2, time_start = $3, time_end = $4 WHERE id = $5", QueuesTable)
	_, err := s.db.Exec(query, queue.Title, queue.AvailableID, queue.TimeStart, queue.TimeEnd, queue.ID)
	return err
}

func (s *PostgresRepository) DeleteQueue(id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", QueuesTable)
	_, err := s.db.Exec(query, id)
	return err
}
