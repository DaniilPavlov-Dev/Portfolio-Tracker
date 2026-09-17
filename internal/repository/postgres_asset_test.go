package repository

import (
	"PFnPTA/internal/model"
	"context"
	"os"
	"testing"
)

func TestPostgresAssetRepositoryCreateAndFindByID(t *testing.T) {
	ctx := context.Background()

	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		t.Fatal("DATABASE_URL is not set")
	}

	repo, err := NewPostgresRepository(ctx, connString)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = repo.Close(ctx)
	})

	assetRepo := NewPostgresAssetRepository(repo)

	asset := &model.Asset{
		Symbol: "TEST_BTC",
		Name:   "Test Bitcoin",
	}

	err = assetRepo.Create(ctx, asset)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := assetRepo.Delete(ctx, asset.ID); err != nil {
			t.Errorf("failed to cleanup asset: %v", err)
		}
	})

	foundAsset, err := assetRepo.FindByID(ctx, asset.ID)
	if err != nil {
		t.Fatal(err)
	}

	if foundAsset.ID != asset.ID {
		t.Fatalf("expected ID %d, got %d", asset.ID, foundAsset.ID)
	}

	if foundAsset.Symbol != asset.Symbol {
		t.Fatalf("expected symbol %s, got %s", asset.Symbol, foundAsset.Symbol)
	}

	if foundAsset.Name != asset.Name {
		t.Fatalf("expected name %s, got %s", asset.Name, foundAsset.Name)
	}
}

func TestPostgresAssetRepositoryFindAll(t *testing.T) {
	ctx := context.Background()

	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		t.Fatal("DATABASE_URL is not set")
	}

	repo, err := NewPostgresRepository(ctx, connString)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = repo.Close(ctx)
	})

	assetRepo := NewPostgresAssetRepository(repo)

	asset := &model.Asset{
		Symbol: "TEST_ETH",
		Name:   "Test Ethereum",
	}

	err = assetRepo.Create(ctx, asset)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := assetRepo.Delete(ctx, asset.ID); err != nil {
			t.Errorf("failed to cleanup asset: %v", err)
		}
	})

	assets, err := assetRepo.FindAll(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if assets == nil {
		t.Fatal("expected assets, got nil")
	}

	var found bool

	for _, item := range assets {
		if item.ID == asset.ID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("asset with ID %d was not found", asset.ID)
	}
}
