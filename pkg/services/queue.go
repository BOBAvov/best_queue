package services

import (
	"sso/models"
)

func (s *AuthService) CreateQueue(queue models.Queue) (int, error) {
	return s.repo.CreateQueue(queue)
}

func (s *AuthService) GetQueueByID(id int) (models.Queue, error) {
	return s.repo.GetQueueByID(id)
}

func (s *AuthService) GetAllQueues() ([]models.Queue, error) {
	return s.repo.GetAllQueues()
}

func (s *AuthService) UpdateQueue(queue models.Queue) error {
	return s.repo.UpdateQueue(queue)
}

func (s *AuthService) DeleteQueue(id int) error {
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
