package repository

import (
	"PFnPTA/internal/model"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type PostgresTransactionRepository struct {
	repo *PostgresRepository
}

var _ TransactionRepository = (*PostgresTransactionRepository)(nil)

func NewPostgresTransactionRepository(repo *PostgresRepository) *PostgresTransactionRepository {
	return &PostgresTransactionRepository{repo: repo}
}

func (r *PostgresTransactionRepository) FindByID(ctx context.Context, id int64) (*model.Transaction, error) {
	var transaction model.Transaction

	err := r.repo.db.QueryRow(
		ctx,
		`SELECT id, user_id, asset_id, type, quantity, price, commission, created_at
		FROM transactions
		WHERE id = $1`,
		id,
	).Scan(
		&transaction.ID,
		&transaction.UserID,
		&transaction.AssetID,
		&transaction.Type,
		&transaction.Quantity,
		&transaction.Price,
		&transaction.Commission,
		&transaction.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTransactionNotFound
		}

		return nil, err
	}

	return &transaction, nil
}

func (r *PostgresTransactionRepository) FindByIDAndUserID(ctx context.Context, transactionID int64, userID int64) (*model.Transaction, error) {
	var transaction model.Transaction

	err := r.repo.db.QueryRow(
		ctx,
		`SELECT id, user_id, asset_id, type, quantity, price, commission, crated_at
		FROM transactions
		WHERE id = $1 AND user_id = $2`,
		transactionID,
		userID,
	).Scan(
		&transaction.ID,
		&transaction.UserID,
		&transaction.AssetID,
		&transaction.Type,
		&transaction.Quantity,
		&transaction.Price,
		&transaction.Commission,
		&transaction.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTransactionNotFound
		}

		return nil, err
	}

	return &transaction, nil
}

func (r *PostgresTransactionRepository) FindByUserID(ctx context.Context, userID int64) ([]*model.Transaction, error) {
	rows, err := r.repo.db.Query(
		ctx,
		`SELECT id, user_id, asset_id, type, quantity, price, commission, created_at
		FROM transactions
		WHERE user_id = $1
		ORDER BY id`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*model.Transaction

	for rows.Next() {
		var transaction model.Transaction

		if err := rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.AssetID,
			&transaction.Type,
			&transaction.Quantity,
			&transaction.Price,
			&transaction.Commission,
			&transaction.CreatedAt,
		); err != nil {
			return nil, err
		}

		transactions = append(transactions, &transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *PostgresTransactionRepository) Create(ctx context.Context, transaction *model.Transaction) error {
	return r.repo.db.QueryRow(
		ctx,
		`INSERT INTO transactions (
			user_id,
			asset_id,
			type,
			quantity,
			price,
			commission
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`,
		transaction.UserID,
		transaction.AssetID,
		transaction.Type,
		transaction.Quantity,
		transaction.Price,
		transaction.Commission,
	).Scan(
		&transaction.ID,
		&transaction.CreatedAt,
	)
}
