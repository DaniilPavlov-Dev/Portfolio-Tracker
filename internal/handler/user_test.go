package handler

import (
	"PFnPTA/internal/repository"
	"PFnPTA/internal/service"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var jwtService = service.NewJWTService("test-secret")

func TestUserHandler_Register(t *testing.T) {
	repo := repository.NewMemoryUserRepository()
	userService := service.NewUserService(repo)
	handler := NewUserHandler(userService, jwtService)

	body := bytes.NewBufferString(`{
		"email": "test@example.com",
		"password": "password123"
	}`)

	req := httptest.NewRequest(
		http.MethodPost,
		"/register",
		body,
	)

	req = req.WithContext(context.Background())

	rec := httptest.NewRecorder()

	handler.Register(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusCreated)
	}

	var response struct {
		ID        int64     `json:"id"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"createdAt"`
	}

	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Email != "test@example.com" {
		t.Errorf("got email %q, want %q", response.Email, "test@example.com")
	}

	if bytes.Contains(rec.Body.Bytes(), []byte("PasswordHash")) {
		t.Error("response contains PasswordHash")
	}

	if !bytes.Contains(rec.Body.Bytes(), []byte(`"id"`)) {
		t.Error("response does not contain id field")
	}

	if !bytes.Contains(rec.Body.Bytes(), []byte(`"email"`)) {
		t.Error("response does not contain email field")
	}
}

func TestUserHandler_Register_InvalidBody(t *testing.T) {
	repo := repository.NewMemoryUserRepository()
	service := service.NewUserService(repo)
	handler := NewUserHandler(service, jwtService)

	req := httptest.NewRequest(
		http.MethodPost,
		"/register",
		bytes.NewBufferString(`invalid json`),
	)

	rec := httptest.NewRecorder()

	handler.Register(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestUserHandler_Register_MethodNotAllowed(t *testing.T) {
	repo := repository.NewMemoryUserRepository()
	service := service.NewUserService(repo)
	handler := NewUserHandler(service, jwtService)

	req := httptest.NewRequest(
		http.MethodGet,
		"/register",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.Register(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func TestUserHandler_Register_ServiceError(t *testing.T) {
	repo := repository.NewMemoryUserRepository()
	service := service.NewUserService(repo)
	handler := NewUserHandler(service, jwtService)

	body := bytes.NewBufferString(`{
		"email": "invalid-email",
		"password": "password123"
	}`)

	req := httptest.NewRequest(
		http.MethodPost,
		"/register",
		body,
	)

	rec := httptest.NewRecorder()

	handler.Register(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
