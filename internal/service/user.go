package service

import (
	"PFnPTA/internal/model"
	"PFnPTA/internal/repository"
	"context"
	"errors"
	"net/mail"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepository
}

func (s *UserService) Register(ctx context.Context, email, password string) (*model.User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, errors.New("invalid email")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}
	if utf8.RuneCountInString(password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}
	_, err := s.repo.FindByEmail(ctx, email)
	if err == nil {
		return nil, errors.New("email already exists")
	}
	if errors.Is(err, repository.ErrUserNotFound) {
		passwordHash, err := bcrypt.GenerateFromPassword(
			[]byte(password),
			bcrypt.DefaultCost,
		)
		if err != nil {
			return nil, err
		}
		user := &model.User{Email: email, PasswordHash: string(passwordHash)}
		err = s.repo.Create(ctx, user)
		if err != nil {
			if errors.Is(err, repository.ErrUserAlreadyExists) {
				return nil, errors.New("email already exists")
			}

			return nil, err
		}
		return user, nil
	}
	return nil, err
}

func (s *UserService) FindByID(ctx context.Context, id int64) (*model.User, error) {
	return s.repo.FindByID(ctx, id)
}

func NewUserService(r repository.UserRepository) *UserService {
	return &UserService{
		repo: r,
	}
}

func (s *UserService) Login(ctx context.Context, email, password string) (*model.User, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)

	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}
