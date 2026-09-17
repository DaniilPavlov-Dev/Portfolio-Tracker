package repository

import (
	"PFnPTA/internal/model"
	"context"
)

type AssetRepository interface {
	FindByID(ctx context.Context, id int64) (*model.Asset, error)
	FindAll(ctx context.Context) ([]*model.Asset, error)
	Create(ctx context.Context, asset *model.Asset) error
	Delete(ctx context.Context, id int64) error
}
