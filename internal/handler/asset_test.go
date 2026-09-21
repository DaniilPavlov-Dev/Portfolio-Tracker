package handler

import (
	"PFnPTA/internal/model"
	"PFnPTA/internal/repository"
	"PFnPTA/internal/service"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestAssetHandler_GetAll(t *testing.T) {
	repo := repository.NewMemoryAssetRepository()

	asset1 := &model.Asset{
		Symbol: "BTC",
		Name:   "Bitcoin",
	}

	asset2 := &model.Asset{
		Symbol: "ETH",
		Name:   "Ethereum",
	}

	if err := repo.Create(context.Background(), asset1); err != nil {
		t.Fatalf("failed to create asset1: %v", err)
	}

	if err := repo.Create(context.Background(), asset2); err != nil {
		t.Fatalf("failed to create asset2: %v", err)
	}

	assetService := service.NewAssetService(repo)
	handler := NewAssetHandler(assetService)

	req := httptest.NewRequest(
		http.MethodGet,
		"/assets",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.GetAll(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	var assets []*model.Asset

	if err := json.Unmarshal(rec.Body.Bytes(), &assets); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(assets) != 2 {
		t.Fatalf("got %d assets, want %d", len(assets), 2)
	}
}

func TestAssetHandler_GetAll_MethodNotAllowed(t *testing.T) {
	repo := repository.NewMemoryAssetRepository()
	assetService := service.NewAssetService(repo)
	handler := NewAssetHandler(assetService)

	req := httptest.NewRequest(
		http.MethodPost,
		"/assets",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.GetAll(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func TestAssetHandler_GetAll_Empty(t *testing.T) {
	repo := repository.NewMemoryAssetRepository()
	assetService := service.NewAssetService(repo)
	handler := NewAssetHandler(assetService)

	req := httptest.NewRequest(
		http.MethodGet,
		"/assets",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.GetAll(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	var assets []*model.Asset

	if err := json.Unmarshal(rec.Body.Bytes(), &assets); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(assets) != 0 {
		t.Fatalf("got %d assets, want 0", len(assets))
	}
}

func TestAssetHandler_GetByID(t *testing.T) {
	repo := repository.NewMemoryAssetRepository()

	asset := &model.Asset{
		Symbol: "BTC",
		Name:   "Bitcoin",
	}

	if err := repo.Create(context.Background(), asset); err != nil {
		t.Fatalf("failed to create asset: %v", err)
	}

	assetService := service.NewAssetService(repo)
	handler := NewAssetHandler(assetService)

	req := httptest.NewRequest(
		http.MethodGet,
		"/assets/"+strconv.FormatInt(asset.ID, 10),
		nil,
	)

	rec := httptest.NewRecorder()

	handler.GetByID(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	var response model.Asset

	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.ID != asset.ID {
		t.Errorf("got ID %d, want %d", response.ID, asset.ID)
	}

	if response.Symbol != asset.Symbol {
		t.Errorf("got symbol %q, want %q", response.Symbol, asset.Symbol)
	}

	if response.Name != asset.Name {
		t.Errorf("got name %q, want %q", response.Name, asset.Name)
	}
}

func TestAssetHandler_GetByID_NotFound(t *testing.T) {
	repo := repository.NewMemoryAssetRepository()
	assetService := service.NewAssetService(repo)
	handler := NewAssetHandler(assetService)

	req := httptest.NewRequest(
		http.MethodGet,
		"/assets/999",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.GetByID(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestAssetHandler_GetByID_InvalidID(t *testing.T) {
	repo := repository.NewMemoryAssetRepository()
	assetService := service.NewAssetService(repo)
	handler := NewAssetHandler(assetService)

	req := httptest.NewRequest(
		http.MethodGet,
		"/assets/abc",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.GetByID(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestAssetHandler_GetByID_MethodNotAllowed(t *testing.T) {
	repo := repository.NewMemoryAssetRepository()
	assetService := service.NewAssetService(repo)
	handler := NewAssetHandler(assetService)

	req := httptest.NewRequest(
		http.MethodPost,
		"/asset/1",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.GetByID(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}
