package service

import (
	"PFnPTA/internal/model"
	"testing"
	"time"
)

func TestToUserResponse(t *testing.T) {
	createAt := time.Now()

	user := &model.User{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: "secret-hash",
		CreatedAt:    createAt,
	}

	response := ToUserResponse(user)
	if response == nil {
		t.Fatal("expected response, got nil")
	}
	if response.ID != user.ID {
		t.Errorf("got ID %d, want %d", response.ID, user.ID)
	}

	if response.Email != user.Email {
		t.Errorf("got Email %q, want %q", response.Email, user.Email)
	}

	if response.CreatedAt != user.CreatedAt {
		t.Errorf("got CreatedAt %v, want %v", response.CreatedAt, user.CreatedAt)
	}
}
