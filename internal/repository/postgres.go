package repository

import (
	"PFnPTA/internal/model"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type PostgresRepository struct {
	db *pgx.Conn
}

var _ UserRepository = (*PostgresRepository)(nil)

func NewPostgresRepository(ctx context.Context, connString string) (*PostgresRepository, error) {
	db, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{db: db}, nil
}

func (r *PostgresRepository) Close(ctx context.Context) error {
	return r.db.Close(ctx)
}

func (r *PostgresRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	err := r.db.QueryRow(
		ctx,
		`SELECT id, email, password_hash, created_at 
		FROM users WHERE 
		email=$1`,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil

}

func (r *PostgresRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.QueryRow(
		ctx,
		`INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, created_at`,
		user.Email,
		user.PasswordHash,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)
}

func (r *PostgresRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User

	err := r.db.QueryRow(
		ctx,
		`SELECT id, email, password_hash, created_at
		FROM users
		WHERE id = $1`,
		id,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *PostgresRepository) FindAssetByID(ctx context.Context, id int64) (*model.Asset, error) {
	var asset model.Asset

	err := r.db.QueryRow(
		ctx,
		`SELECT id, symbol, name, created_at
		FROM assets
		WHERE id = $1`,
		id,
	).Scan(
		&asset.ID,
		&asset.Symbol,
		&asset.Name,
		&asset.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAssetNotFound
		}
		return nil, err
	}

	return &asset, nil
}

func (r *PostgresRepository) FindAllAssets(ctx context.Context) ([]*model.Asset, error) {
	rows, err := r.db.Query(
		ctx,
		`SELECT id, symbol, name, created_at
		FROM assets
		ORDER BY id`,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []*model.Asset

	for rows.Next() {
		var asset model.Asset

		if err := rows.Scan(
			&asset.ID,
			&asset.Symbol,
			&asset.Name,
			&asset.CreatedAt,
		); err != nil {
			return nil, err
		}

		assets = append(assets, &asset)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return assets, nil
}
