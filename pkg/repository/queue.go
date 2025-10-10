package repository

import (
	"fmt"
	"sso/models"
)

func (r *PostgresRepository) CreateQueue(queue models.Queue) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (title, time_start, time_end) VALUES ($1, $2, $3) RETURNING id", QueuesTable)
	err := r.db.QueryRow(query, queue.Title, queue.TimeStart, queue.TimeEnd).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *PostgresRepository) GetQueueByID(id int) (models.Queue, error) {
	var queue models.Queue
	query := fmt.Sprintf("SELECT id, title, time_start, time_end FROM %s WHERE id = $1", QueuesTable)
	err := r.db.Get(&queue, query, id)
	if err != nil {
		return queue, err
	}
	return queue, nil
}

func (r *PostgresRepository) GetAllQueues() ([]models.Queue, error) {
	var queues []models.Queue
	query := fmt.Sprintf("SELECT id, title, time_start, time_end FROM %s ORDER BY time_start DESC", QueuesTable)
	err := r.db.Select(&queues, query)
	if err != nil {
		return nil, err
	}
	return queues, nil
}

func (r *PostgresRepository) UpdateQueue(queue models.Queue) error {
	query := fmt.Sprintf("UPDATE %s SET title = $1, time_start = $2, time_end = $3 WHERE id = $4", QueuesTable)
	_, err := r.db.Exec(query, queue.Title, queue.TimeStart, queue.TimeEnd, queue.ID)
	return err
}

func (r *PostgresRepository) DeleteQueue(id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", QueuesTable)
	_, err := r.db.Exec(query, id)
	return err
}
