package handler

import (
	"PFnPTA/internal/middleware"
	"PFnPTA/internal/model"
	"PFnPTA/internal/repository"
	"PFnPTA/internal/service"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTransactionAPI_Create(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()
	transactionService := service.NewTransactionService(repo, assetRepo)
	handler := NewTransactionHandler(transactionService)

	jwtService := service.NewJWTService("test-secret")
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	asset := &model.Asset{
		Symbol: "BTC",
		Name:   "Bitcoin",
	}
	if err := assetRepo.Create(context.Background(), asset); err != nil {
		t.Fatal(err)
	}

	token, err := jwtService.GenerateToken(1)
	if err != nil {
		t.Fatal(err)
	}

	requestBody := fmt.Sprintf(`{
		"assetId": %d,
		"type": "BUY",
		"quantity": 0.5,
		"price": 60000,
		"commission": 10
	}`, asset.ID)

	req := httptest.NewRequest(
		http.MethodPost,
		"/transactions",
		strings.NewReader(requestBody),
	)

	req.Header.Set(
		"Authorization",
		"Bearer "+token,
	)

	rec := httptest.NewRecorder()

	authMiddleware.RequireAuth(http.HandlerFunc(handler.Create)).ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusCreated)
	}

	transactions, err := repo.FindByUserID(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}

	if len(transactions) != 1 {
		t.Fatalf("got %d transactions, want 1", len(transactions))
	}

	transaction := transactions[0]

	if transaction.UserID != 1 {
		t.Errorf("got user ID %d, want %d", transaction.UserID, 1)
	}

	if transaction.AssetID != asset.ID {
		t.Errorf("got asset ID %d, want %d", transaction.AssetID, asset.ID)
	}

	if transaction.Type != model.TransactionBuy {
		t.Errorf("got type %v, want %v", transaction.Type, model.TransactionBuy)
	}

	if transaction.Quantity != 0.5 {
		t.Errorf("got quantity %f, want %f", transaction.Quantity, 0.5)
	}
}

func TestTransactionAPI_GetByUserID(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()

	transaction1 := &model.Transaction{
		UserID:   1,
		AssetID:  1,
		Type:     model.TransactionBuy,
		Quantity: 0.5,
		Price:    60000,
	}

	transaction2 := &model.Transaction{
		UserID:   2,
		AssetID:  1,
		Type:     model.TransactionBuy,
		Quantity: 1,
		Price:    70000,
	}

	if err := repo.Create(context.Background(), transaction1); err != nil {
		t.Fatal(err)
	}

	if err := repo.Create(context.Background(), transaction2); err != nil {
		t.Fatal(err)
	}

	transactionService := service.NewTransactionService(repo, assetRepo)
	handler := NewTransactionHandler(transactionService)

	jwtService := service.NewJWTService("test-secret")
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	token, err := jwtService.GenerateToken(1)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/transactions",
		nil,
	)

	req.Header.Set(
		"Authorization",
		"Bearer "+token,
	)

	rec := httptest.NewRecorder()

	authMiddleware.RequireAuth(
		http.HandlerFunc(handler.GetByUserID),
	).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"got status %d, want %d",
			rec.Code,
			http.StatusOK,
		)
	}

	var response []model.Transaction

	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response) != 1 {
		t.Fatalf("got %d transactions, want 1", len(response))
	}

	if response[0].UserID != 1 {
		t.Errorf("got user ID %d, want %d", response[0].UserID, 1)
	}
}
