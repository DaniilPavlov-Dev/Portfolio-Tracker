package service

import (
	"PFnPTA/internal/model"
	"PFnPTA/internal/repository"
	"context"
	"errors"
	"testing"
)

func TestTransactionService_Create_InvalidUserID(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()
	service := NewTransactionService(repo, assetRepo)

	transaction := &model.Transaction{
		UserID: 0,
	}

	err := service.Create(context.Background(), transaction)

	if err == nil || err.Error() != "user id must be greater than 0" {
		t.Errorf("got error %v, want %v", err, "user id must be greater than 0")
	}
}

func TestTransactionService_Create_NegativeUserID(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()
	service := NewTransactionService(repo, assetRepo)

	transaction := &model.Transaction{
		UserID: -1,
	}

	err := service.Create(context.Background(), transaction)

	if err == nil || err.Error() != "user id must be greater than 0" {
		t.Errorf("got error %v, want %v", err, "user id must be greater than 0")
	}
}

func TestTransactionService_Create_InvalidAssetID(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()
	service := NewTransactionService(repo, assetRepo)

	transaction := &model.Transaction{
		UserID:  1,
		AssetID: 0,
	}

	err := service.Create(context.Background(), transaction)

	if err == nil || err.Error() != "asset id must be greater than 0" {
		t.Errorf("got error %v, want %v", err, "asset id must be greater than 0")
	}
}

func TestTransactionService_Create_NegativeAssetID(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()
	service := NewTransactionService(repo, assetRepo)

	transaction := &model.Transaction{
		UserID:  1,
		AssetID: -1,
	}

	err := service.Create(context.Background(), transaction)

	if err == nil || err.Error() != "asset id must be greater than 0" {
		t.Errorf("got error %v, want %v", err, "asset id must be greater than 0")
	}
}

func TestTransactionService_Create_InvalidType(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()
	service := NewTransactionService(repo, assetRepo)

	transaction := &model.Transaction{
		UserID:  1,
		AssetID: 1,
		Type:    "HOLD",
	}

	err := service.Create(context.Background(), transaction)

	if err == nil || err.Error() != "invalid transaction type" {
		t.Errorf("got error %v, want %v", err, "invalid transaction type")
	}
}

func TestTransactionService_Create_EmptyType(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()
	service := NewTransactionService(repo, assetRepo)

	transaction := &model.Transaction{
		UserID:  1,
		AssetID: 1,
		Type:    "",
	}

	err := service.Create(context.Background(), transaction)

	if err == nil || err.Error() != "invalid transaction type" {
		t.Errorf("got error %v, want %v", err, "invalid transaction type")
	}
}

func TestTransactionService_Create_InvalidQuantity(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()
	service := NewTransactionService(repo, assetRepo)

	transaction := &model.Transaction{
		UserID:   1,
		AssetID:  1,
		Type:     model.TransactionBuy,
		Quantity: 0,
	}

	err := service.Create(context.Background(), transaction)

	if err == nil || err.Error() != "quantity must be greater than 0" {
		t.Errorf("got error %v, want %v", err, "quantity must be greater than 0")
	}
}

func TestTransactionService_Create_NegativeQuantity(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()
	service := NewTransactionService(repo, assetRepo)

	transaction := &model.Transaction{
		UserID:   1,
		AssetID:  1,
		Type:     model.TransactionBuy,
		Quantity: -1,
	}

	err := service.Create(context.Background(), transaction)

	if err == nil || err.Error() != "quantity must be greater than 0" {
		t.Errorf("got error %v, want %v", err, "quantity must be greater than 0")
	}
}

func TestTransactionService_Create_InvalidPrice(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()
	service := NewTransactionService(repo, assetRepo)

	transaction := &model.Transaction{
		UserID:   1,
		AssetID:  1,
		Type:     model.TransactionBuy,
		Quantity: 1,
		Price:    0,
	}

	err := service.Create(context.Background(), transaction)

	if err == nil || err.Error() != "price must be greater than 0" {
		t.Errorf("got error %v, want %v", err, "price must be greater than 0")
	}
}

func TestTransactionService_Create_NegativePrice(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()
	service := NewTransactionService(repo, assetRepo)

	transaction := &model.Transaction{
		UserID:   1,
		AssetID:  1,
		Type:     model.TransactionBuy,
		Quantity: 1,
		Price:    -1,
	}

	err := service.Create(context.Background(), transaction)

	if err == nil || err.Error() != "price must be greater than 0" {
		t.Errorf("got error %v, want %v", err, "price must be greater than 0")
	}
}

func TestTransactionService_Create_NegativeCommission(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()
	service := NewTransactionService(repo, assetRepo)

	transaction := &model.Transaction{
		UserID:     1,
		AssetID:    1,
		Type:       model.TransactionBuy,
		Quantity:   1,
		Price:      1,
		Commission: -5,
	}

	err := service.Create(context.Background(), transaction)

	if err == nil || err.Error() != "commission cannot be negative" {
		t.Errorf("got error %v, want %v", err, "commission cannot be negative")
	}
}

func TestTransactionService_Create_Success(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()
	service := NewTransactionService(repo, assetRepo)

	asset := &model.Asset{
		Symbol: "BTC",
		Name:   "Bitcoin",
	}

	if err := assetRepo.Create(context.Background(), asset); err != nil {
		t.Fatal(err)
	}

	transaction := &model.Transaction{
		UserID:     1,
		AssetID:    asset.ID,
		Type:       model.TransactionBuy,
		Quantity:   0.5,
		Price:      60000,
		Commission: 10,
	}

	err := service.Create(context.Background(), transaction)

	if err != nil {
		t.Fatalf("failed to create transaction: %v", err)
	}

	if transaction.ID < 1 {
		t.Errorf("got transaction id %d, want >= 1", transaction.ID)
	}

	got, err := repo.FindByID(context.Background(), transaction.ID)
	if err != nil {
		t.Errorf("got error %v, want %v", err, repository.ErrTransactionNotFound)
	}
	if got.UserID != transaction.UserID {
		t.Errorf("got user id %d, want %d", got.UserID, transaction.UserID)
	}
	if got.AssetID != transaction.AssetID {
		t.Errorf("got asset id %d, want %d", got.AssetID, transaction.AssetID)
	}
	if got.Type != transaction.Type {
		t.Errorf("got transaction type id %s, want %s", got.Type, transaction.Type)
	}
	if got.Quantity != transaction.Quantity {
		t.Errorf("got quantity %f, want %f", got.Quantity, transaction.Quantity)
	}
	if got.Price != transaction.Price {
		t.Errorf("got price %f, want %f", got.Price, transaction.Price)
	}
}

func TestTransactionService_FindByID(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()
	service := NewTransactionService(repo, assetRepo)

	asset := &model.Asset{
		Symbol: "BTC",
		Name:   "Bitcoin",
	}

	if err := assetRepo.Create(context.Background(), asset); err != nil {
		t.Fatal(err)
	}

	transaction := &model.Transaction{
		UserID:     1,
		AssetID:    asset.ID,
		Type:       model.TransactionBuy,
		Quantity:   0.5,
		Price:      60000,
		Commission: 10,
	}

	err := service.Create(context.Background(), transaction)

	if err != nil {
		t.Fatalf("failed to create transaction: %v", err)
	}

	if transaction.ID < 1 {
		t.Errorf("got transaction id %d, want >= 1", transaction.ID)
	}

	got, err := service.FindByID(context.Background(), transaction.ID)
	if err != nil {
		t.Errorf("got error %v, want %v", err, repository.ErrTransactionNotFound)
	}
	if got.ID != transaction.ID {
		t.Errorf("expected id %d, got %d", transaction.ID, got.ID)
	}
}

func TestTransactionService_FindByID_NotFound(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()
	service := NewTransactionService(repo, assetRepo)

	_, err := service.FindByID(context.Background(), 999)

	if !errors.Is(err, repository.ErrTransactionNotFound) {
		t.Errorf("got %v, want %v", err, repository.ErrTransactionNotFound)
	}
}

func TestTransactionService_FindByUserID(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()
	service := NewTransactionService(repo, assetRepo)

	asset := &model.Asset{
		Symbol: "BTC",
		Name:   "Bitcoin",
	}

	if err := assetRepo.Create(context.Background(), asset); err != nil {
		t.Fatal(err)
	}

	transaction1 := &model.Transaction{
		UserID:     1,
		AssetID:    asset.ID,
		Type:       model.TransactionBuy,
		Quantity:   0.5,
		Price:      60000,
		Commission: 10,
	}
	transaction2 := &model.Transaction{
		UserID:     1,
		AssetID:    asset.ID,
		Type:       model.TransactionBuy,
		Quantity:   0.5,
		Price:      60000,
		Commission: 10,
	}
	transaction3 := &model.Transaction{
		UserID:     2,
		AssetID:    asset.ID,
		Type:       model.TransactionBuy,
		Quantity:   0.5,
		Price:      60000,
		Commission: 10,
	}

	if err := service.Create(context.Background(), transaction1); err != nil {
		t.Fatalf("failed to create transaction: %v", err)
	}

	if err := service.Create(context.Background(), transaction2); err != nil {
		t.Fatalf("failed to create transaction: %v", err)
	}

	if err := service.Create(context.Background(), transaction3); err != nil {
		t.Fatalf("failed to create transaction: %v", err)
	}

	got, err := service.FindByUserID(context.Background(), 1)
	if err != nil {
		t.Fatalf("failed to find transactions: %v", err)
	}

	if len(got) != 2 {
		t.Errorf("got %d transactions, want 2", len(got))
	}
	for _, transaction := range got {
		if transaction.UserID != 1 {
			t.Errorf("got user %d, want 1", transaction.UserID)
		}
	}
}

func TestTransactionService_FindByUserID_Empty(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()
	service := NewTransactionService(repo, assetRepo)

	got, err := service.FindByUserID(context.Background(), 999)

	if err != nil {
		t.Fatalf("failed to find transactions: %v", err)
	}

	if len(got) != 0 {
		t.Errorf("got %d transactions, want 0", len(got))
	}
}

func TestTransactionService_FindByIDAndUserID(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()
	service := NewTransactionService(repo, assetRepo)

	asset := &model.Asset{
		Symbol: "BTC",
		Name:   "Bitcoin",
	}

	if err := assetRepo.Create(context.Background(), asset); err != nil {
		t.Fatal(err)
	}

	transaction := &model.Transaction{
		UserID:     1,
		AssetID:    asset.ID,
		Type:       model.TransactionBuy,
		Quantity:   0.5,
		Price:      60000,
		Commission: 10,
	}

	if err := service.Create(context.Background(), transaction); err != nil {
		t.Fatalf("failed to create transaction: %v", err)
	}

	got, err := service.FindByIDAndUserID(
		context.Background(),
		transaction.ID,
		1,
	)
	if err != nil {
		t.Fatalf("got error %v", err)
	}

	if got.ID != transaction.ID {
		t.Errorf("expected id %d, got %d", transaction.ID, got.ID)
	}

	if got.UserID != 1 {
		t.Errorf("expected user id %d, got %d", 1, got.UserID)
	}
}

func TestTransactionService_FindByIDAndUserID_NotOwner(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()
	service := NewTransactionService(repo, assetRepo)

	asset := &model.Asset{
		Symbol: "BTC",
		Name:   "Bitcoin",
	}

	if err := assetRepo.Create(context.Background(), asset); err != nil {
		t.Fatal(err)
	}

	transaction := &model.Transaction{
		UserID:     2,
		AssetID:    asset.ID,
		Type:       model.TransactionBuy,
		Quantity:   0.5,
		Price:      60000,
		Commission: 10,
	}

	if err := service.Create(context.Background(), transaction); err != nil {
		t.Fatalf("failed to create transaction: %v", err)
	}

	_, err := service.FindByIDAndUserID(
		context.Background(),
		transaction.ID,
		1,
	)

	if !errors.Is(err, repository.ErrTransactionNotFound) {
		t.Fatalf(
			"got error %v, want %v",
			err,
			repository.ErrTransactionNotFound,
		)
	}
}
