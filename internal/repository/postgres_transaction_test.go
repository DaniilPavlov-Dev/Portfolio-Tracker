package repository

import (
	"PFnPTA/internal/model"
	"context"
	"errors"
	"os"
	"testing"
)

func TestPostgresTransactionRepositoryCreate(t *testing.T) {
	ctx := context.Background()

	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		t.Fatal("DATABASE_URL is not set")
	}

	repo, err := NewPostgresRepository(ctx, connString)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_ = repo.Close(ctx)
	})

	transactionRepo := NewPostgresTransactionRepository(repo)

	transaction := &model.Transaction{
		UserID:     6,
		AssetID:    10,
		Type:       model.TransactionBuy,
		Quantity:   0.5,
		Price:      100000,
		Commission: 10,
	}

	err = transactionRepo.Create(ctx, transaction)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_, err := repo.db.Exec(
			ctx,
			`DELETE FROM transactions WHERE id = $1`,
			transaction.ID,
		)
		if err != nil {
			t.Errorf("failed to cleanup transactions: %v", err)
		}
	})

	if transaction.ID == 0 {
		t.Fatal("expected transaction ID to be set")
	}

	if transaction.CreatedAt.IsZero() {
		t.Fatal("expected CreatedAt to be set")
	}

	foundTransaction, err := transactionRepo.FindByID(ctx, transaction.ID)
	if err != nil {
		t.Fatal(err)
	}

	if foundTransaction.ID != transaction.ID {
		t.Fatalf("expected ID %d, got %d", transaction.ID, foundTransaction.ID)

	}

	if foundTransaction.UserID != transaction.UserID {
		t.Fatalf("expected UserID %d, got %d", transaction.UserID, foundTransaction.UserID)
	}

	if foundTransaction.AssetID != transaction.AssetID {
		t.Fatalf("expected AssetID %d, got %d", transaction.AssetID, foundTransaction.AssetID)
	}

	if foundTransaction.Type != transaction.Type {
		t.Fatalf("expected Type %s, got %s", transaction.Type, foundTransaction.Type)
	}

	if foundTransaction.Quantity != transaction.Quantity {
		t.Fatalf("expected Quantity %f, got %f", transaction.Quantity, foundTransaction.Quantity)
	}

	if foundTransaction.Price != transaction.Price {
		t.Fatalf("expected Price %f, got %f", transaction.Price, foundTransaction.Price)
	}

	if foundTransaction.Commission != transaction.Commission {
		t.Fatalf("expected Commission %f, got %f", transaction.Commission, foundTransaction.Commission)
	}
}

func TestPostgresTransactionRepositoryFindByIDNotFound(t *testing.T) {
	ctx := context.Background()

	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		t.Fatal("DATABASE_URL is not set")
	}

	repo, err := NewPostgresRepository(ctx, connString)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_ = repo.Close(ctx)
	})

	transactionRepo := NewPostgresTransactionRepository(repo)

	_, err = transactionRepo.FindByID(ctx, 9999999)
	if !errors.Is(err, ErrTransactionNotFound) {
		t.Fatalf("expected ErrTransactionNotFound, got %v", err)
	}
}

func TestPostgresTransactionRepositoryFindByUserID(t *testing.T) {
	ctx := context.Background()

	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		t.Fatal("DATABASE_URL is not set")
	}

	repo, err := NewPostgresRepository(ctx, connString)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_ = repo.Close(ctx)
	})

	transactionRepo := NewPostgresTransactionRepository(repo)

	transaction := &model.Transaction{
		UserID:     6,
		AssetID:    10,
		Type:       model.TransactionBuy,
		Quantity:   0.5,
		Price:      100000,
		Commission: 10,
	}

	err = transactionRepo.Create(ctx, transaction)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_, err := repo.db.Exec(
			ctx,
			`DELETE FROM transactions WHERE id = $1`,
			transaction.ID,
		)
		if err != nil {
			t.Errorf("failed to cleanup transactions: %v", err)
		}
	})

	transactions, err := transactionRepo.FindByUserID(ctx, transaction.UserID)
	if err != nil {
		t.Fatal(err)
	}

	if len(transactions) == 0 {
		t.Fatal("expected at least one transaction")
	}

	var found bool

	for _, item := range transactions {
		if item.ID == transaction.ID {
			found = true

			if item.UserID != transaction.UserID {
				t.Fatalf("expected UserID %d, got %d", transaction.UserID, item.UserID)
			}

			if item.AssetID != transaction.AssetID {
				t.Fatalf("expected AssetID %d, got %d", transaction.AssetID, item.AssetID)
			}

			if item.Type != transaction.Type {
				t.Fatalf("expected Type %s, got %s", transaction.Type, item.Type)
			}

			if item.Quantity != transaction.Quantity {
				t.Fatalf("expected Quantity %f, got %f", transaction.Quantity, item.Quantity)
			}

			if item.Price != transaction.Price {
				t.Fatalf("expected Price %f, got %f", transaction.Price, item.Price)
			}

			if item.Commission != transaction.Commission {
				t.Fatalf("expected Commission %f, got %f", transaction.Commission, item.Commission)
			}

			break
		}
	}

	if !found {
		t.Fatalf("transaction with ID %d, was not found", transaction.ID)
	}
}
