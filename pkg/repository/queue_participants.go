package repository

import (
	"fmt"
	"sso/models"
)

// JoinQueue добавляет пользователя в очередь
func (r *PostgresRepository) JoinQueue(queueID, userID int) (int, error) {
	// Проверяем, не находится ли пользователь уже в очереди
	var existingID int
	checkQuery := fmt.Sprintf("SELECT id FROM %s WHERE queue_id = $1 AND user_id = $2 AND is_active = true", QueueParticipantsTable)
	err := r.db.QueryRow(checkQuery, queueID, userID).Scan(&existingID)
	if err == nil {
		return 0, fmt.Errorf("user is already in queue")
	}

	// Получаем следующую позицию в очереди
	position, err := r.GetNextQueuePosition(queueID)
	if err != nil {
		return 0, err
	}

	// Добавляем пользователя в очередь
	var id int
	query := fmt.Sprintf("INSERT INTO %s (queue_id, user_id, position) VALUES ($1, $2, $3) RETURNING id", QueueParticipantsTable)
	err = r.db.QueryRow(query, queueID, userID, position).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// LeaveQueue удаляет пользователя из очереди
func (r *PostgresRepository) LeaveQueue(queueID, userID int) error {
	query := fmt.Sprintf("UPDATE %s SET is_active = false WHERE queue_id = $1 AND user_id = $2", QueueParticipantsTable)
	result, err := r.db.Exec(query, queueID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found in queue")
	}

	return nil
}

// GetQueueParticipants возвращает всех участников очереди
func (r *PostgresRepository) GetQueueParticipants(queueID int) ([]models.QueueParticipant, error) {
	var participants []models.QueueParticipant
	query := fmt.Sprintf("SELECT id, queue_id, user_id, position, joined_at, is_active FROM %s WHERE queue_id = $1 AND is_active = true ORDER BY position", QueueParticipantsTable)
	err := r.db.Select(&participants, query, queueID)
	if err != nil {
		return nil, err
	}

	return participants, nil
}

// GetUserQueuePosition возвращает позицию пользователя в очереди
func (r *PostgresRepository) GetUserQueuePosition(queueID, userID int) (int, error) {
	var position int
	query := fmt.Sprintf("SELECT position FROM %s WHERE queue_id = $1 AND user_id = $2 AND is_active = true", QueueParticipantsTable)
	err := r.db.QueryRow(query, queueID, userID).Scan(&position)
	if err != nil {
		return 0, err
	}

	return position, nil
}

// ShiftQueue удаляет первого пользователя из очереди и сдвигает остальных
func (r *PostgresRepository) ShiftQueue(queueID int) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Удаляем первого пользователя из очереди
	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE queue_id = $1 AND position = 1 AND is_active = true", QueueParticipantsTable)
	_, err = tx.Exec(deleteQuery, queueID)
	if err != nil {
		return err
	}

	// Сдвигаем позиции остальных пользователей
	updateQuery := fmt.Sprintf("UPDATE %s SET position = position - 1 WHERE queue_id = $1 AND is_active = true", QueueParticipantsTable)
	_, err = tx.Exec(updateQuery, queueID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetNextQueuePosition возвращает следующую позицию в очереди
func (r *PostgresRepository) GetNextQueuePosition(queueID int) (int, error) {
	var position int
	query := fmt.Sprintf("SELECT COALESCE(MAX(position), 0) + 1 FROM %s WHERE queue_id = $1 AND is_active = true", QueueParticipantsTable)
	err := r.db.QueryRow(query, queueID).Scan(&position)
	if err != nil {
		return 0, err
	}

	return position, nil
}
