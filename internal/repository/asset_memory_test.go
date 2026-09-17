package repository

import (
	"PFnPTA/internal/model"
	"context"
	"testing"
)

func TestMemoryAssetRepository_FindByID(t *testing.T) {
	repo := NewMemoryAssetRepository()

	asset := &model.Asset{
		Symbol: "BTC",
		Name:   "Bitcoin",
	}

	if err := repo.Create(context.Background(), asset); err != nil {
		t.Fatalf("failed to create asset: %v", err)
	}

	got, err := repo.FindByID(context.Background(), asset.ID)
	if err != nil {
		t.Fatalf("failed to find asset: %v", err)
	}

	if got.ID != asset.ID {
		t.Errorf("got ID %d, want %d", got.ID, asset.ID)
	}

	if got.Symbol != asset.Symbol {
		t.Errorf("got symbol %q, want %q", got.Symbol, asset.Symbol)
	}

	if got.Name != asset.Name {
		t.Errorf("got name %q, want %q", got.Name, asset.Name)
	}
}

func TestMemoryAssetRepository_FindByID_NotFound(t *testing.T) {
	repo := NewMemoryAssetRepository()

	_, err := repo.FindByID(context.Background(), 999)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err != ErrAssetNotFound {
		t.Errorf("got error %v, want %v", err, ErrAssetNotFound)
	}
}

func TestMemoryAssetRepository_Create(t *testing.T) {
	repo := NewMemoryAssetRepository()

	asset := &model.Asset{
		Symbol: "ETH",
		Name:   "Ethereum",
	}

	if err := repo.Create(context.Background(), asset); err != nil {
		t.Fatalf("failed to create asset: %v", err)
	}

	if asset.ID == 0 {
		t.Error("expected asset ID to be set")
	}

	got, err := repo.FindByID(context.Background(), asset.ID)
	if err != nil {
		t.Fatalf("failed to find created asset: %v", err)
	}

	if got != asset {
		t.Error("repository returned a different asset")
	}
}

func TestMemoryAssetRepository_FindAll(t *testing.T) {
	repo := NewMemoryAssetRepository()

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

	assets, err := repo.FindAll(context.Background())
	if err != nil {
		t.Fatalf("failed to find all assets: %v", err)
	}

	if len(assets) != 2 {
		t.Fatalf("got %d assets, want %d", len(assets), 2)
	}
}

func TestMemoryAssetRepository_Delete(t *testing.T) {
	repo := NewMemoryAssetRepository()

	asset := &model.Asset{
		Symbol: "BTC",
		Name:   "Bitcoin",
	}

	if err := repo.Create(context.Background(), asset); err != nil {
		t.Fatalf("failed to create asset: %v", err)
	}

	if err := repo.Delete(context.Background(), asset.ID); err != nil {
		t.Fatalf("failed to delete asset: %v", err)
	}

	_, err := repo.FindByID(context.Background(), asset.ID)
	if err != ErrAssetNotFound {
		t.Errorf("got error %v, want %v", err, ErrAssetNotFound)
	}
}

func TestMemoryAssetRepository_Create_SetsCreatedAt(t *testing.T) {
	repo := NewMemoryAssetRepository()

	asset := &model.Asset{
		Symbol: "BTC",
		Name:   "Bitcoin",
	}

	if err := repo.Create(context.Background(), asset); err != nil {
		t.Fatalf("failed to create asset: %v", err)
	}

	if asset.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}
