package repository

import (
	"PFnPTA/internal/model"
	"context"
	"errors"
	"time"
)

var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyExists = errors.New("user already exists")

type MemoryUserRepository struct {
	users map[string]*model.User
}

var _ UserRepository = (*MemoryUserRepository)(nil)

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users: make(map[string]*model.User),
	}
}

func (r *MemoryUserRepository) FindByEmail(
	_ context.Context,
	email string,
) (*model.User, error) {
	user, ok := r.users[email]
	if !ok {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (r *MemoryUserRepository) Create(
	_ context.Context,
	user *model.User,
) error {
	if _, ok := r.users[user.Email]; ok {
		return ErrUserAlreadyExists
	}
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}
	r.users[user.Email] = user
	return nil
}

func (r *MemoryUserRepository) FindByID(_ context.Context, id int64) (*model.User, error) {
	for _, user := range r.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}
