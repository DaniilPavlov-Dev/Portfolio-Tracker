package repository

import (
	"PFnPTA/internal/model"
	"context"
	"os"
	"testing"
)

func TestNewPostgresRepository(t *testing.T) {
	ctx := context.Background()

	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		t.Fatal("DATABASE_URL is not set")
	}

	repo, err := NewPostgresRepository(ctx, connString)
	if err != nil {
		t.Fatalf("failed to coonect to database: %v", err)
	}

	defer repo.Close(ctx)
}

func TestPostgresRepositoryCreateAndFindByEmail(t *testing.T) {
	ctx := context.Background()

	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		t.Fatal("DATABASE_URL is not set")
	}

	repo, err := NewPostgresRepository(ctx, connString)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer repo.Close(ctx)

	user := &model.User{
		Email:        "postgres-test@example.com",
		PasswordHash: "test-hash",
	}

	err = repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	defer func() {
		_, _ = repo.db.Exec(
			ctx,
			"DELETE FROM users WHERE email = $1",
			user.Email,
		)
	}()

	got, err := repo.FindByEmail(ctx, user.Email)
	if err != nil {
		t.Fatalf("failed to find user: %v", err)
	}
	if got.Email != user.Email {
		t.Errorf("expected email %q, got %q", got.Email, user.Email)
	}
	if got.PasswordHash != user.PasswordHash {
		t.Errorf("expected password hash %q, got %q", got.PasswordHash, user.PasswordHash)
	}
	if got.ID == 0 {
		t.Error("expected user ID to be generated")
	}
	if got.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be generated")
	}
}
