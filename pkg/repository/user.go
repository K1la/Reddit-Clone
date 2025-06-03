package repository

import (
	"errors"
	"redditclone/pkg/models"
	"sync"
)

type (
	InMemoryUserRepo struct {
		users map[string]*models.User
		mu    sync.RWMutex
	}

	UserRepo interface {
		Create(userName, hashPassword string) (*models.User, error)
		GetByUsername(username string) (*models.User, error)
	}
)

func NewInMemoryUserRepo() *InMemoryUserRepo {
	return &InMemoryUserRepo{
		users: make(map[string]*models.User),
	}
}

func (r *InMemoryUserRepo) Create(userName, hashPassword string) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exist := r.users[userName]; exist {
		return nil, errors.New("username already exists")
	}
	user := &models.User{
		Username: userName,
		Password: hashPassword,
	}
	r.users[userName] = user
	return user, nil
}

func (r *InMemoryUserRepo) GetByUsername(username string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exist := r.users[username]
	if !exist {
		return nil, errors.New("user not found")
	}
	return user, nil
}
