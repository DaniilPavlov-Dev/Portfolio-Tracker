package repository

import (
	"PFnPTA/internal/model"
	"context"
)

type TransactionRepository interface {
	FindByID(ctx context.Context, id int64) (*model.Transaction, error)
	FindByUserID(ctx context.Context, id int64) ([]*model.Transaction, error)
	Create(ctx context.Context, transaction *model.Transaction) error
}
