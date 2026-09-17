package repository

import (
	"PFnPTA/internal/model"
	"context"
	"testing"
)

func TestMemoryTransactionRepository_FindByID(t *testing.T) {
	repo := NewMemoryTransactionRepository()

	transaction := &model.Transaction{
		ID:         1,
		UserID:     10,
		AssetID:    5,
		Type:       model.TransactionBuy,
		Quantity:   0.5,
		Price:      60000,
		Commission: 10,
	}

	if err := repo.Create(context.Background(), transaction); err != nil {
		t.Fatalf("failed to create transaction: %v", err)
	}

	got, err := repo.FindByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("failed to find transaction: %v", err)
	}

	if got.ID != 1 {
		t.Errorf("got ID %d, want %d", got.ID, transaction.ID)
	}
}

func TestMemoryTransactionRepository_FindByID_NotFound(t *testing.T) {
	repo := NewMemoryTransactionRepository()

	_, err := repo.FindByID(context.Background(), 999)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err != ErrTransactionNotFound {
		t.Errorf("got error %v, want %v", err, ErrTransactionNotFound)
	}
}

func TestMemoryTransactionRepository_FindByUserID(t *testing.T) {
	repo := NewMemoryTransactionRepository()

	transaction1 := &model.Transaction{
		UserID:     10,
		AssetID:    5,
		Type:       model.TransactionBuy,
		Quantity:   0.5,
		Price:      60000,
		Commission: 10,
	}
	transaction2 := &model.Transaction{
		UserID:     10,
		AssetID:    5,
		Type:       model.TransactionSell,
		Quantity:   0.5,
		Price:      60000,
		Commission: 10,
	}
	transaction3 := &model.Transaction{
		UserID:     20,
		AssetID:    5,
		Type:       model.TransactionBuy,
		Quantity:   0.5,
		Price:      60000,
		Commission: 10,
	}

	if err := repo.Create(context.Background(), transaction1); err != nil {
		t.Fatalf("failed to create transaction: %v", err)
	}
	if err := repo.Create(context.Background(), transaction2); err != nil {
		t.Fatalf("failed to create transaction: %v", err)
	}
	if err := repo.Create(context.Background(), transaction3); err != nil {
		t.Fatalf("failed to create transaction: %v", err)
	}

	got, err := repo.FindByUserID(context.Background(), 10)
	if err != nil {
		t.Fatalf("failed to find transactions: %v", err)
	}

	if len(got) != 2 {
		t.Errorf("got %d transactions, want %d", len(got), 2)
	}

	for _, transaction := range got {
		if transaction.UserID != 10 {
			t.Errorf("got transaction with UserID %d, wabt %d", transaction.ID, 10)
		}
	}
}

func TestMemoryTransactionRepository_FindByUserID_Empty(t *testing.T) {
	repo := NewMemoryTransactionRepository()

	got, err := repo.FindByUserID(context.Background(), 999)

	if err != nil {
		t.Fatalf("got error: %v", err)
	}

	if len(got) != 0 {
		t.Errorf("got %d transactions, want %d", len(got), 0)
	}
}
