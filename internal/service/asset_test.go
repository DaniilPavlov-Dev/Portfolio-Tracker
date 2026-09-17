package service

import (
	"PFnPTA/internal/model"
	"PFnPTA/internal/repository"
	"context"
	"testing"
)

type fakeAssetRepository struct {
	assets map[int64]*model.Asset
}

func (r *fakeAssetRepository) FindByID(_ context.Context, id int64) (*model.Asset, error) {
	asset, ok := r.assets[id]
	if !ok {
		return nil, repository.ErrAssetNotFound
	}
	return asset, nil
}

func (r *fakeAssetRepository) FindAll(_ context.Context) ([]*model.Asset, error) {
	assets := make([]*model.Asset, 0, len(r.assets))
	for _, asset := range r.assets {
		assets = append(assets, asset)
	}
	return assets, nil
}

func (r *fakeAssetRepository) Create(_ context.Context, asset *model.Asset) error {
	r.assets[asset.ID] = asset
	return nil
}

func (r *fakeAssetRepository) Delete(_ context.Context, id int64) error {
	if _, ok := r.assets[id]; !ok {
		return repository.ErrAssetNotFound
	}
	delete(r.assets, id)
	return nil
}

var _ repository.AssetRepository = (*fakeAssetRepository)(nil)

func TestAssetService_GetByID(t *testing.T) {
	asset := &model.Asset{
		ID:     1,
		Symbol: "BTC",
		Name:   "Bitcoin",
	}

	repo := &fakeAssetRepository{
		assets: map[int64]*model.Asset{
			asset.ID: asset,
		},
	}

	service := NewAssetService(repo)

	got, err := service.GetByID(context.Background(), asset.ID)
	if err != nil {
		t.Fatalf("got error: %v", err)
	}

	if got != asset {
		t.Error("service returned a different asset")
	}
}

func TestAssetService_GetByID_NotFound(t *testing.T) {
	repo := &fakeAssetRepository{
		assets: make(map[int64]*model.Asset),
	}

	service := NewAssetService(repo)

	_, err := service.GetByID(context.Background(), 999)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err != repository.ErrAssetNotFound {
		t.Errorf("got error %v, want %v", err, repository.ErrAssetNotFound)
	}
}

func TestAssetService_GetAll(t *testing.T) {
	asset1 := &model.Asset{
		ID:     1,
		Symbol: "BTC",
		Name:   "Bitcoin",
	}

	asset2 := &model.Asset{
		ID:     2,
		Symbol: "ETH",
		Name:   "Ethereum",
	}

	repo := &fakeAssetRepository{
		assets: map[int64]*model.Asset{
			asset1.ID: asset1,
			asset2.ID: asset2,
		},
	}

	service := NewAssetService(repo)

	assets, err := service.GetAll(context.Background())
	if err != nil {
		t.Fatalf("got error: %v", err)
	}

	if len(assets) != 2 {
		t.Fatalf("got %d assets, want %d", len(assets), 2)
	}
}

func TestAssetService_Create(t *testing.T) {
	repo := &fakeAssetRepository{
		assets: make(map[int64]*model.Asset),
	}

	service := NewAssetService(repo)

	asset := &model.Asset{
		ID:     1,
		Symbol: "BTC",
		Name:   "Bitcoin",
	}

	err := service.Create(context.Background(), asset)
	if err != nil {
		t.Fatalf("got error: %v", err)
	}

	saved, ok := repo.assets[asset.ID]
	if !ok {
		t.Fatal("asset was not saved")
	}

	if saved != asset {
		t.Errorf("saved asset is different")
	}
}

func TestAssetService_Delete(t *testing.T) {
	asset := &model.Asset{
		ID:     1,
		Symbol: "BTC",
		Name:   "Bitcoin",
	}

	repo := &fakeAssetRepository{
		assets: map[int64]*model.Asset{
			asset.ID: asset,
		},
	}

	service := NewAssetService(repo)

	err := service.Delete(context.Background(), asset.ID)
	if err != nil {
		t.Fatalf("got error: %v", err)
	}

	if _, ok := repo.assets[asset.ID]; ok {
		t.Error("asset was not deleted")
	}
}

func TestAssetService_Create_EmptySymbol(t *testing.T) {
	asset := &model.Asset{
		ID:     1,
		Symbol: "",
		Name:   "Bitcoin",
	}

	repo := &fakeAssetRepository{
		assets: map[int64]*model.Asset{
			asset.ID: asset,
		},
	}

	service := NewAssetService(repo)

	err := service.Create(context.Background(), asset)

	if err == nil || err.Error() != "symbol is required" {
		t.Errorf("got error: %v, want %v", err, "symbol is required")
	}
}

func TestAssetService_Create_EmptyName(t *testing.T) {
	asset := &model.Asset{
		ID:     1,
		Symbol: "BTC",
		Name:   "",
	}

	repo := &fakeAssetRepository{
		assets: map[int64]*model.Asset{
			asset.ID: asset,
		},
	}

	service := NewAssetService(repo)

	err := service.Create(context.Background(), asset)

	if err == nil || err.Error() != "name is required" {
		t.Errorf("got error: %v, want %v", err, "name is required")
	}
}

func TestAssetService_Create_EmptySymbolAndName(t *testing.T) {
	asset := &model.Asset{
		ID:     1,
		Symbol: "",
		Name:   "",
	}

	repo := &fakeAssetRepository{
		assets: map[int64]*model.Asset{
			asset.ID: asset,
		},
	}

	service := NewAssetService(repo)

	err := service.Create(context.Background(), asset)

	if err == nil || err.Error() != "symbol is required; name is required" {
		t.Errorf("got error: %v, want %v", err, "symbol is required; name is required")
	}
}
