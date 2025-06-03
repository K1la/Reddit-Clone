package repository

import (
	"errors"
	"redditclone/pkg/models"

	"github.com/google/uuid"
)

type (
	InMemorySessionRepo struct {
		sessions map[string]*models.Session
	}

	SessionRepo interface {
		Create(userName string) (*models.Session, error)
	}
)

func NewInMemorySessionRepo() *InMemorySessionRepo {
	return &InMemorySessionRepo{
		sessions: make(map[string]*models.Session),
	}
}

func (r *InMemorySessionRepo) Create(userName string) (*models.Session, error) {

	if _, exist := r.sessions[userName]; exist {
		return nil, errors.New("username already exists")
	}
	session := &models.Session{
		ID:       uuid.NewString(),
		Username: userName,
	}
	r.sessions[userName] = session
	return session, nil
}
