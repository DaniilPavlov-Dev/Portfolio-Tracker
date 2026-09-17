package service

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateToken(t *testing.T) {
	service := NewJWTService("test-secret")

	tokenString, err := service.GenerateToken(42)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	if tokenString == "" {
		t.Fatal("expected token, got empty string")
	}

	token, err := jwt.Parse(
		tokenString,
		func(t *jwt.Token) (any, error) {
			return []byte("test-secret"), nil
		},
	)
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}
	if !token.Valid {
		t.Fatal("expected token tyo be valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("expected MapClaims")
	}

	if claims["user_id"] != float64(42) {
		t.Errorf("expected user_id 42, got %v", claims["user_id"])
	}
}

func TestParseToken(t *testing.T) {
	service := NewJWTService("test-secret")

	tokenString, err := service.GenerateToken(42)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	userID, err := service.ParseToken(tokenString)
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}

	if userID != 42 {
		t.Errorf("expected suer ID 42, got %d", userID)
	}
}

func TestParseTokenInvalid(t *testing.T) {
	service := NewJWTService("test-secret")

	_, err := service.ParseToken("invalid-token")
	if err == nil {
		t.Fatalf("expected error for invalid token")
	}
}
