package repository

import (
	"PFnPTA/internal/model"
	"context"
	"testing"
	"time"
)

func TestMemoryUserRepository_FindByEmail(t *testing.T) {
	user := &model.User{
		ID:    1,
		Email: "test@example.com",
	}
	repo := NewMemoryUserRepository()
	repo.users[user.Email] = user
	got, err := repo.FindByEmail(context.Background(), user.Email)
	if err != nil {
		t.Fatalf("Test is failing, err: %q", err)
	}
	if got.ID != user.ID || got.Email != user.Email {
		t.Errorf("User with ID %v not found", got.ID)
	}
}

func TestMemoryUserRepository_FindByEmail_NotFound(t *testing.T) {
	repo := NewMemoryUserRepository()
	_, err := repo.FindByEmail(context.Background(), "unknown@example.com")
	if err == nil {
		t.Errorf("Expected error, not nil")
	}
}

func TestMemoryUserRepository_Create(t *testing.T) {
	repo := NewMemoryUserRepository()

	user := &model.User{
		ID:    1,
		Email: "new@example.com",
	}

	err := repo.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	got, err := repo.FindByEmail(context.Background(), user.Email)
	if err != nil {
		t.Fatalf("FindByEmail failed: %v", err)
	}
	if got.ID != user.ID || got.Email != user.Email {
		t.Errorf("got user %+v, want %+v", got, user)
	}
}

func TestMemoryUserRepository_Create_SetsCreatedAt(t *testing.T) {
	repo := NewMemoryUserRepository()

	user := &model.User{
		ID:    1,
		Email: "test@example.com",
	}

	err := repo.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if user.CreatedAt.IsZero() {
		t.Error("CreatedAt was not set")
	}
}

func TestMemoryUserRepository_Create_PreservesCreatedAt(t *testing.T) {
	repo := NewMemoryUserRepository()
	createdAt := time.Date(
		2025,
		time.January,
		1,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	user := &model.User{
		ID:        1,
		Email:     "test@example.com",
		CreatedAt: createdAt,
	}

	err := repo.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if !user.CreatedAt.Equal(createdAt) {
		t.Errorf("CreatedAt was changed: got %v, want %v", user.CreatedAt, createdAt)
	}
}

func TestMemoryUserRepository_Create_DuplicateEmail(t *testing.T) {
	repo := NewMemoryUserRepository()

	user1 := &model.User{
		ID:    1,
		Email: "test@example.com",
	}

	err := repo.Create(context.Background(), user1)
	if err != nil {
		t.Fatalf("first Create failed: %v", err)
	}

	user2 := &model.User{
		ID:    2,
		Email: "test@example.com",
	}

	err = repo.Create(context.Background(), user2)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "user already exists" {
		t.Errorf("got error %q, want %q", err.Error(), "user already exists")
	}
}
