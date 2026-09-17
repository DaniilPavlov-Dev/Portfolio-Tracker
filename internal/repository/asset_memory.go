package repository

import (
	"PFnPTA/internal/model"
	"context"
	"errors"
	"time"
)

var ErrAssetNotFound = errors.New("asset not found")

type MemoryAssetRepository struct {
	assets map[int64]*model.Asset
	nextID int64
}

func NewMemoryAssetRepository() *MemoryAssetRepository {
	return &MemoryAssetRepository{
		assets: make(map[int64]*model.Asset),
		nextID: 1,
	}
}

func (r *MemoryAssetRepository) FindByID(_ context.Context, id int64) (*model.Asset, error) {
	asset, ok := r.assets[id]
	if !ok {
		return nil, ErrAssetNotFound
	}

	return asset, nil
}

func (r *MemoryAssetRepository) FindAll(_ context.Context) ([]*model.Asset, error) {
	assets := make([]*model.Asset, 0, len(r.assets))

	for _, asset := range r.assets {
		assets = append(assets, asset)
	}

	return assets, nil
}

func (r *MemoryAssetRepository) Create(_ context.Context, asset *model.Asset) error {
	asset.ID = r.nextID
	if asset.CreatedAt.IsZero() {
		asset.CreatedAt = time.Now()
	}
	r.nextID++

	r.assets[asset.ID] = asset

	return nil
}

func (r *MemoryAssetRepository) Delete(_ context.Context, id int64) error {
	if _, ok := r.assets[id]; !ok {
		return ErrAssetNotFound
	}

	delete(r.assets, id)

	return nil
}
