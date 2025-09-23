package repository

import (
	"errors"
	"sync"
)

type Authorization interface {
	CreateUser(user User) (int, error)
	GetUserByUsername(username string) (User, error)
}

type AuthMemoryRepo struct {
	mu      sync.Mutex
	users   map[string]User
	counter int
}

func NewAuthMemoryRepo() *AuthMemoryRepo {
	return &AuthMemoryRepo{
		users:   make(map[string]User),
		counter: 0,
	}
}

func (r *AuthMemoryRepo) CreateUser(user User) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.Tg_name]; exists {
		return 0, errors.New("user already exists")
	}

	user.Id = r.counter
	r.users[user.Tg_name] = user
	r.counter++

	return user.Id, nil
}

func (r *AuthMemoryRepo) GetUserByUsername(username string) (User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if user, exists := r.users[username]; exists {
		return user, nil
	}

	return User{}, errors.New("user not found")
}
