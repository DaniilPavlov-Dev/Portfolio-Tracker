package service

import (
	"PFnPTA/internal/model"
	"PFnPTA/internal/repository"
	"context"
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

type fakeUserRepository struct {
	users     map[string]*model.User
	createErr error
}

func (r *fakeUserRepository) FindByEmail(_ context.Context, email string) (*model.User, error) {
	user, ok := r.users[email]
	if !ok {
		return nil, repository.ErrUserNotFound
	}
	return user, nil
}

func (r *fakeUserRepository) FindByID(_ context.Context, id int64) (*model.User, error) {
	for _, user := range r.users {
		if user.ID == 0 {
			return user, nil
		}
	}

	return nil, repository.ErrUserNotFound
}

func (r *fakeUserRepository) Create(_ context.Context, user *model.User) error {
	if r.createErr != nil {
		return r.createErr
	}
	r.users[user.Email] = user
	return nil
}

var _ repository.UserRepository = (*fakeUserRepository)(nil)

func TestRegister(t *testing.T) {
	repo := &fakeUserRepository{
		users: make(map[string]*model.User),
	}
	service := NewUserService(repo)

	_, err := service.Register(
		context.Background(),
		"test@example.com",
		"password123",
	)

	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	savedUser, err := repo.FindByEmail(
		context.Background(),
		"test@example.com",
	)

	if err != nil {
		t.Fatalf("FindByEmail failed: %v", err)
	}
	if savedUser.Email != "test@example.com" {
		t.Errorf("got email %q, want %q", savedUser.Email, "test@example.com")
	}
	if savedUser.PasswordHash == "password123" {
		t.Error("password was stored without hashing")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(savedUser.PasswordHash),
		[]byte("password123"),
	)
	if err != nil {
		t.Errorf("password hash verification failed: %v", err)
	}
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	repo := &fakeUserRepository{
		users: map[string]*model.User{
			"test@example.com": {
				ID:    1,
				Email: "test@example.com",
			},
		},
	}
	service := NewUserService(repo)

	_, err := service.Register(
		context.Background(),
		"test@example.com",
		"password123",
	)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err.Error() != "email already exists" {
		t.Errorf("got error %q, want %q", err.Error(), "email already exists")
	}
}

func TestRegister_CreateError(t *testing.T) {
	createErr := errors.New("create failed")
	repo := &fakeUserRepository{
		users:     make(map[string]*model.User),
		createErr: createErr,
	}
	service := NewUserService(repo)
	_, err := service.Register(
		context.Background(),
		"test@example.com",
		"password123",
	)

	if !errors.Is(err, createErr) {
		t.Errorf("got %q, want %q", err, createErr)
	}
}

func TestRegister_WrongPassword(t *testing.T) {
	repo := &fakeUserRepository{
		users: make(map[string]*model.User),
	}
	service := NewUserService(repo)

	user, err := service.Register(
		context.Background(),
		"test@example.com",
		"password123",
	)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte("wrong-password"),
	)

	if err == nil {
		t.Errorf("wrong password was accepted")
	}

}

func TestRegister_PasswordValidation(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		wantErr     bool
		wantErrText string
	}{
		{name: "empty password", password: "", wantErr: true, wantErrText: "password is required"},
		{name: "short password", password: "1234567", wantErr: true, wantErrText: "password must be at least 8 characters"},
		{name: "minimum length password", password: "12345678", wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeUserRepository{
				users: make(map[string]*model.User),
			}
			service := NewUserService(repo)

			_, err := service.Register(
				context.Background(),
				"test@example.com",
				tt.password,
			)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}

				if err.Error() != tt.wantErrText {
					t.Errorf("got error %q, want %q", err.Error(), tt.wantErrText)
				}

				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestRegister_EmailValidation(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		wantErr     bool
		wantErrText string
	}{
		{name: "empty email", email: "", wantErr: true, wantErrText: "email is required"},
		{name: "invalid email", email: "test", wantErr: true, wantErrText: "invalid email"},
		{name: "valid email", email: "test@example.com", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeUserRepository{
				users: make(map[string]*model.User),
			}

			service := NewUserService(repo)
			_, err := service.Register(
				context.Background(),
				tt.email,
				"password123",
			)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if err.Error() != tt.wantErrText {
					t.Errorf("got error %q, want %q", err.Error(), tt.wantErrText)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
