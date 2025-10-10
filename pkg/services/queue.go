package services

import (
	"sso/models"
	"time"
)

func (s *AuthService) CreateQueue(queue models.Queue) (int, error) {
	queue.TimeStart = time.Now()
	if queue.TimeAdd == 0 {
		queue.TimeAdd = 5 // по умолчанию 5 часов
	}
	queue.TimeEnd = time.Now().Add(time.Hour * time.Duration(queue.TimeAdd))
	return s.repo.CreateQueue(queue)
}

func (s *AuthService) GetQueueById(id string) (models.Queue, error) {
	return s.repo.GetQueueById(id)
}

func (s *AuthService) GetAllQueues() ([]models.Queue, error) {
	return s.repo.GetAllQueues()
}

func (s *AuthService) UpdateQueue(queue models.Queue) error {
	return s.repo.UpdateQueue(queue)
}

func (s *AuthService) DeleteQueue(id string) error {
	return s.repo.DeleteQueue(id)
}

// Queue Participants methods
func (s *AuthService) JoinQueue(queueID, userID int) (int, error) {
	return s.repo.JoinQueue(queueID, userID)
}

func (s *AuthService) LeaveQueue(queueID, userID int) error {
	return s.repo.LeaveQueue(queueID, userID)
}

func (s *AuthService) GetQueueParticipants(queueID int) ([]models.QueueParticipant, error) {
	return s.repo.GetQueueParticipants(queueID)
}

func (s *AuthService) ShiftQueue(queueID int) error {
	return s.repo.ShiftQueue(queueID)
}

// User management methods
func (s *AuthService) GetUserByID(id int) (models.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *AuthService) GetAllUsers() ([]models.User, error) {
	return s.repo.GetAllUsers()
}

func (s *AuthService) UpdateUser(id int, user models.User) error {
	return s.repo.UpdateUser(id, user)
}

func (s *AuthService) DeleteUser(id int) error {
	return s.repo.DeleteUser(id)
}

func (s *AuthService) GetGroupIdByCode(code string) (int, error) {
	return s.repo.GetGroupIdByCode(code)
}
