package handler

import (
	"PFnPTA/internal/middleware"
	"PFnPTA/internal/model"
	"PFnPTA/internal/repository"
	"PFnPTA/internal/service"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestTransactionHandler_Create(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()

	asset := &model.Asset{
		Symbol: "BTC",
		Name:   "Bitcoin",
	}

	if err := assetRepo.Create(context.Background(), asset); err != nil {
		t.Fatal(err)
	}

	service := service.NewTransactionService(repo, assetRepo)
	handler := NewTransactionHandler(service)

	request := struct {
		AssetID    int64   `json:"assetId"`
		Type       string  `json:"type"`
		Quantity   float64 `json:"quantity"`
		Price      float64 `json:"price"`
		Commission float64 `json:"commission"`
	}{
		AssetID:    asset.ID,
		Type:       "BUY",
		Quantity:   0.5,
		Price:      60000,
		Commission: 10,
	}

	bodyJSON, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}

	body := bytes.NewBuffer(bodyJSON)

	req := httptest.NewRequest(
		http.MethodPost,
		"/transactions",
		body,
	)

	req = req.WithContext(middleware.ContextWithUserID(req.Context(), 1))

	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusCreated)
	}

	var response model.Transaction

	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.ID == 0 {
		t.Fatalf("expected transaction ID to be set")
	}

	if response.UserID != 1 {
		t.Errorf("got user ID %d, want %d", response.UserID, 1)
	}

	if response.AssetID != 1 {
		t.Errorf("got asset ID %d, want %d", response.AssetID, 1)
	}

	if response.Type != model.TransactionBuy {
		t.Errorf("got type %v, want %v", response.Type, model.TransactionBuy)
	}

	if response.Quantity != 0.5 {
		t.Errorf("got quantity %f, want %f", response.Quantity, 0.5)
	}

	if response.Price != 60000 {
		t.Errorf("got price %f, want %f", response.Price, 60000.0)
	}

	if response.Commission != 10 {
		t.Errorf("got commission %f, want %f", response.Commission, 10.0)
	}

	if response.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestTransactionHandler_Create_InvalidJSON(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()
	service := service.NewTransactionService(repo, assetRepo)
	handler := NewTransactionHandler(service)

	body := bytes.NewBufferString(`invalid json`)

	req := httptest.NewRequest(
		http.MethodPost,
		"/transactions",
		body,
	)

	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestTransactionHandler_Create_InvalidTransaction(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()
	service := service.NewTransactionService(repo, assetRepo)
	handler := NewTransactionHandler(service)

	body := bytes.NewBufferString(`{
		"assetId": 1,
		"type": "BUY",
		"quantity": 0,
		"price": 60000,
		"commission": 10
	}`)

	req := httptest.NewRequest(
		http.MethodPost,
		"/transactions",
		body,
	)
	req = req.WithContext(middleware.ContextWithUserID(req.Context(), 1))

	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestTransactionHandler_Create_MethodNotAllowed(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()
	service := service.NewTransactionService(repo, assetRepo)
	handler := NewTransactionHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/transactions",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func TestTransactionHandler_GetByID(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()

	transaction := &model.Transaction{
		UserID:     1,
		AssetID:    1,
		Type:       model.TransactionBuy,
		Quantity:   0.5,
		Price:      60000,
		Commission: 10,
	}

	if err := repo.Create(context.Background(), transaction); err != nil {
		t.Fatalf("failed to create transaction: %v", err)
	}

	assetRepo := repository.NewMemoryAssetRepository()
	service := service.NewTransactionService(repo, assetRepo)
	handler := NewTransactionHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/transactions/"+strconv.FormatInt(transaction.ID, 10),
		nil,
	)

	req = req.WithContext(
		middleware.ContextWithUserID(req.Context(), 1),
	)

	rec := httptest.NewRecorder()

	handler.GetByID(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	var response model.Transaction

	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.ID == 0 {
		t.Fatalf("expected transaction ID to be set")
	}

	if response.UserID != 1 {
		t.Errorf("got user ID %d, want %d", response.UserID, 1)
	}

	if response.AssetID != 1 {
		t.Errorf("got asset ID %d, want %d", response.AssetID, 1)
	}

	if response.Type != model.TransactionBuy {
		t.Errorf("got type %v, want %v", response.Type, model.TransactionBuy)
	}

	if response.Quantity != 0.5 {
		t.Errorf("got quantity %f, want %f", response.Quantity, 0.5)
	}

	if response.Price != 60000 {
		t.Errorf("got price %f, want %f", response.Price, 60000.0)
	}

	if response.Commission != 10 {
		t.Errorf("got commission %f, want %f", response.Commission, 10.0)
	}

	if response.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestTransactionHandler_GetByID_NotFound(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()
	service := service.NewTransactionService(repo, assetRepo)
	handler := NewTransactionHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/transactions/999",
		nil,
	)

	req = req.WithContext(
		middleware.ContextWithUserID(req.Context(), 1),
	)

	rec := httptest.NewRecorder()

	handler.GetByID(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestTransactionHandler_GetByID_InvalidID(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()
	service := service.NewTransactionService(repo, assetRepo)
	handler := NewTransactionHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/transactions/abc",
		nil,
	)

	req = req.WithContext(
		middleware.ContextWithUserID(req.Context(), 1),
	)

	rec := httptest.NewRecorder()

	handler.GetByID(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestTransactionHandler_GetByID_Forbidden(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()

	transaction := &model.Transaction{
		UserID:     2,
		AssetID:    1,
		Type:       model.TransactionBuy,
		Quantity:   0.5,
		Price:      60000,
		Commission: 10,
	}

	if err := repo.Create(context.Background(), transaction); err != nil {
		t.Fatalf("failed to create transaction: %v", err)
	}

	assetRepo := repository.NewMemoryAssetRepository()
	service := service.NewTransactionService(repo, assetRepo)
	handler := NewTransactionHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/transactions/1",
		nil,
	)

	req = req.WithContext(
		middleware.ContextWithUserID(req.Context(), 1),
	)

	rec := httptest.NewRecorder()

	handler.GetByID(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusNotFound)
	}

}

func TestTransactionHandler_GetByID_MethodNotAllowed(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()
	service := service.NewTransactionService(repo, assetRepo)
	handler := NewTransactionHandler(service)

	req := httptest.NewRequest(
		http.MethodPost,
		"/transactions/1",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.GetByID(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func TestTransactionHandler_GetByUserID(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()

	transaction1 := &model.Transaction{
		UserID:     1,
		AssetID:    1,
		Type:       model.TransactionBuy,
		Quantity:   0.5,
		Price:      60000,
		Commission: 10,
	}
	transaction2 := &model.Transaction{
		UserID:     1,
		AssetID:    1,
		Type:       model.TransactionBuy,
		Quantity:   0.5,
		Price:      60000,
		Commission: 10,
	}
	transaction3 := &model.Transaction{
		UserID:     2,
		AssetID:    1,
		Type:       model.TransactionBuy,
		Quantity:   0.5,
		Price:      60000,
		Commission: 10,
	}

	if err := repo.Create(context.Background(), transaction1); err != nil {
		t.Fatalf("failed to create transaction: %v", err)
	}
	if err := repo.Create(context.Background(), transaction2); err != nil {
		t.Fatalf("failed to create transaction: %v", err)
	}
	if err := repo.Create(context.Background(), transaction3); err != nil {
		t.Fatalf("failed to create transaction: %v", err)
	}

	assetRepo := repository.NewMemoryAssetRepository()
	service := service.NewTransactionService(repo, assetRepo)
	handler := NewTransactionHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/transactions",
		nil,
	)

	req = req.WithContext(middleware.ContextWithUserID(req.Context(), 1))

	rec := httptest.NewRecorder()

	handler.GetByUserID(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	var response []model.Transaction

	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response) != 2 {
		t.Errorf("got %d transactions, want %d", len(response), 2)
	}

	for _, transaction := range response {
		if transaction.UserID != 1 {
			t.Errorf("got user ID %d, want %d", transaction.UserID, 1)
		}
	}
}

func TestTransactionHandler_GetByUserID_Empty(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()
	service := service.NewTransactionService(repo, assetRepo)
	handler := NewTransactionHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/transactions",
		nil,
	)
	req = req.WithContext(middleware.ContextWithUserID(req.Context(), 1))

	rec := httptest.NewRecorder()

	handler.GetByUserID(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	var response []model.Transaction

	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response) != 0 {
		t.Errorf("got %d transactions, want %d", len(response), 0)
	}
}

func TestTransactionHandler_GetByUserID_Unauthorized(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()
	service := service.NewTransactionService(repo, assetRepo)
	handler := NewTransactionHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/transactions?userId",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.GetByUserID(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestTransactionHandler_GetByUserID_MethodNotAllowed(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()
	service := service.NewTransactionService(repo, assetRepo)
	handler := NewTransactionHandler(service)

	req := httptest.NewRequest(
		http.MethodPost,
		"/transactions?userId=1",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.GetByUserID(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func TestTransactionHandler_GetByID_Unauthorized(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()
	service := service.NewTransactionService(repo, assetRepo)
	handler := NewTransactionHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/transactions/1",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.GetByID(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}
