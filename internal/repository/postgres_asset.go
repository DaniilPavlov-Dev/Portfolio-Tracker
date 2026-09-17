package repository

import (
	"PFnPTA/internal/model"
	"context"
)

type PostgresAssetRepository struct {
	repo *PostgresRepository
}

var _ AssetRepository = (*PostgresAssetRepository)(nil)

func NewPostgresAssetRepository(repo *PostgresRepository) *PostgresAssetRepository {
	return &PostgresAssetRepository{
		repo: repo,
	}
}

func (r *PostgresAssetRepository) FindByID(ctx context.Context, id int64) (*model.Asset, error) {
	return r.repo.FindAssetByID(ctx, id)
}

func (r *PostgresAssetRepository) FindAll(ctx context.Context) ([]*model.Asset, error) {
	return r.repo.FindAllAssets(ctx)
}

func (r *PostgresAssetRepository) Create(ctx context.Context, asset *model.Asset) error {
	return r.repo.db.QueryRow(
		ctx,
		`INSERT INTO assets (symbol, name)
		VALUES ($1, $2)
		RETURNING id, created_at`,
		asset.Symbol,
		asset.Name,
	).Scan(
		&asset.ID,
		&asset.CreatedAt,
	)
}

func (r *PostgresAssetRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.repo.db.Exec(
		ctx,
		`DELETE FROM assets
		WHERE id = $1`,
		id,
	)

	if err != nil {
		return err
	}

	return nil
}
