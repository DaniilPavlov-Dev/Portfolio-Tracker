package repository

import (
	"PFnPTA/internal/model"
	"context"
	"errors"
	"sort"
	"time"
)

type MemoryTransactionRepository struct {
	transaction map[int64]*model.Transaction
	nextID      int64
}

var ErrTransactionNotFound = errors.New("transaction not found")

func NewMemoryTransactionRepository() *MemoryTransactionRepository {
	return &MemoryTransactionRepository{
		transaction: make(map[int64]*model.Transaction),
		nextID:      1,
	}
}

func (r *MemoryTransactionRepository) Create(_ context.Context, transaction *model.Transaction) error {
	transaction.ID = r.nextID
	if transaction.CreatedAt.IsZero() {
		transaction.CreatedAt = time.Now()
	}
	r.nextID++

	r.transaction[transaction.ID] = transaction

	return nil
}

func (r *MemoryTransactionRepository) FindByID(_ context.Context, id int64) (*model.Transaction, error) {
	transaction, ok := r.transaction[id]
	if !ok {
		return nil, ErrTransactionNotFound
	}

	return transaction, nil
}

func (r *MemoryTransactionRepository) FindByUserID(_ context.Context, id int64) ([]*model.Transaction, error) {
	transactions := make([]*model.Transaction, 0, len(r.transaction))

	for _, transaction := range r.transaction {
		if id == transaction.UserID {
			transactions = append(transactions, transaction)
		}
	}

	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].ID < transactions[j].ID
	})

	return transactions, nil
}

func (r *MemoryTransactionRepository) FindByIDAndUserID(_ context.Context, transactionID int64, userID int64) (*model.Transaction, error) {
	transaction, ok := r.transaction[transactionID]
	if !ok {
		return nil, ErrTransactionNotFound
	}

	if transaction.UserID != userID {
		return nil, ErrTransactionNotFound
	}

	return transaction, nil
}
