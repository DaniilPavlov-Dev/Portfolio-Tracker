package service

import (
	"PFnPTA/internal/model"
	"PFnPTA/internal/repository"
	"context"
	"errors"
	"strings"
)

type AssetService struct {
	repo repository.AssetRepository
}

func NewAssetService(repo repository.AssetRepository) *AssetService {
	return &AssetService{
		repo: repo,
	}
}

func (s *AssetService) GetByID(ctx context.Context, id int64) (*model.Asset, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *AssetService) GetAll(ctx context.Context) ([]*model.Asset, error) {
	return s.repo.FindAll(ctx)
}

func (s *AssetService) Create(ctx context.Context, asset *model.Asset) error {
	var validationErrors []string

	if asset.Symbol == "" {
		validationErrors = append(validationErrors, "symbol is required")
	}

	if asset.Name == "" {
		validationErrors = append(validationErrors, "name is required")
	}

	if len(validationErrors) > 0 {
		return errors.New(strings.Join(validationErrors, "; "))
	}

	return s.repo.Create(ctx, asset)
}

func (s *AssetService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
