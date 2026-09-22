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

func TestCalculateAvailableQuantity(t *testing.T) {
	transactions := []*model.Transaction{
		{
			AssetID:  1,
			Type:     model.TransactionBuy,
			Quantity: 1,
		},
	}

	quantity := calculateAvailableQuantity(transactions, 1)

	if quantity != 1 {
		t.Fatalf("got quantity %v, want 1", quantity)
	}
}

func TestCalculateAvailableQuantity_BuyAndSell(t *testing.T) {
	transactions := []*model.Transaction{
		{
			AssetID:  1,
			Type:     model.TransactionBuy,
			Quantity: 1,
		},
		{
			AssetID:  1,
			Type:     model.TransactionSell,
			Quantity: 0.4,
		},
	}

	quantity := calculateAvailableQuantity(transactions, 1)

	if quantity != 0.6 {
		t.Fatalf("got quantity %v, want 0.6", quantity)
	}
}

func TestCalculateAvailableQuantity_DifferentAssets(t *testing.T) {
	transactions := []*model.Transaction{
		{
			AssetID:  1,
			Type:     model.TransactionBuy,
			Quantity: 1,
		},
		{
			AssetID:  2,
			Type:     model.TransactionBuy,
			Quantity: 2,
		},
	}

	quantity := calculateAvailableQuantity(transactions, 1)

	if quantity != 1 {
		t.Fatalf("got quantity %v, want 1", quantity)
	}
}

func TestTransactionService_Create_InsufficientQuantity(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()

	asset := &model.Asset{
		Symbol: "BTC",
		Name:   "Bitcoin",
	}

	if err := assetRepo.Create(context.Background(), asset); err != nil {
		t.Fatal(err)
	}

	service := NewTransactionService(repo, assetRepo)

	buy := &model.Transaction{
		UserID:   1,
		AssetID:  asset.ID,
		Type:     model.TransactionBuy,
		Quantity: 1,
		Price:    100,
	}

	if err := service.Create(context.Background(), buy); err != nil {
		t.Fatal(err)
	}

	sell := &model.Transaction{
		UserID:   1,
		AssetID:  asset.ID,
		Type:     model.TransactionSell,
		Quantity: 1.5,
		Price:    120,
	}

	err := service.Create(context.Background(), sell)

	if err == nil {
		t.Fatal("expected insufficient quantity error")
	}

	if err.Error() != "insufficient quantity" {
		t.Fatalf("got error %q, want %q", err.Error(), "insufficient quantity")
	}
}

func TestTransactionService_Create_PartialSell(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()

	asset := &model.Asset{
		Symbol: "BTC",
		Name:   "Bitcoin",
	}

	if err := assetRepo.Create(context.Background(), asset); err != nil {
		t.Fatal(err)
	}

	service := NewTransactionService(repo, assetRepo)

	buy := &model.Transaction{
		UserID:   1,
		AssetID:  asset.ID,
		Type:     model.TransactionBuy,
		Quantity: 1,
		Price:    100,
	}

	if err := service.Create(context.Background(), buy); err != nil {
		t.Fatal(err)
	}

	sell := &model.Transaction{
		UserID:   1,
		AssetID:  asset.ID,
		Type:     model.TransactionSell,
		Quantity: 0.4,
		Price:    120,
	}

	if err := service.Create(context.Background(), sell); err != nil {
		t.Fatal(err)
	}

	transactions, err := repo.FindByUserID(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}

	if len(transactions) != 2 {
		t.Fatalf("expected 2 transaction, got %d", len(transactions))
	}
}

func TestTransactionService_Create_SellAll(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()

	asset := &model.Asset{
		Symbol: "BTC",
		Name:   "Bitcoin",
	}

	if err := assetRepo.Create(context.Background(), asset); err != nil {
		t.Fatal(err)
	}

	service := NewTransactionService(repo, assetRepo)

	buy := &model.Transaction{
		UserID:   1,
		AssetID:  asset.ID,
		Type:     model.TransactionBuy,
		Quantity: 1,
		Price:    100,
	}

	if err := service.Create(context.Background(), buy); err != nil {
		t.Fatal(err)
	}

	sell := &model.Transaction{
		UserID:   1,
		AssetID:  asset.ID,
		Type:     model.TransactionSell,
		Quantity: 1,
		Price:    120,
	}

	if err := service.Create(context.Background(), sell); err != nil {
		t.Fatal(err)
	}

	transactions, err := repo.FindByUserID(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}

	if len(transactions) != 2 {
		t.Fatalf("expected 2 transaction, got %d", len(transactions))
	}
}

func TestTransactionService_Create_SellAllStepByStep(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()

	asset := &model.Asset{
		Symbol: "BTC",
		Name:   "Bitcoin",
	}

	if err := assetRepo.Create(context.Background(), asset); err != nil {
		t.Fatal(err)
	}

	service := NewTransactionService(repo, assetRepo)

	buy := &model.Transaction{
		UserID:   1,
		AssetID:  asset.ID,
		Type:     model.TransactionBuy,
		Quantity: 1,
		Price:    100,
	}

	if err := service.Create(context.Background(), buy); err != nil {
		t.Fatal(err)
	}

	sell1 := &model.Transaction{
		UserID:   1,
		AssetID:  asset.ID,
		Type:     model.TransactionSell,
		Quantity: 0.4,
		Price:    120,
	}

	if err := service.Create(context.Background(), sell1); err != nil {
		t.Fatal(err)
	}

	sell2 := &model.Transaction{
		UserID:   1,
		AssetID:  asset.ID,
		Type:     model.TransactionSell,
		Quantity: 0.6,
		Price:    120,
	}

	if err := service.Create(context.Background(), sell2); err != nil {
		t.Fatal(err)
	}

	transactions, err := repo.FindByUserID(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}

	if len(transactions) != 3 {
		t.Fatalf("expected 3 transaction, got %d", len(transactions))
	}
}

func TestCalculateAvailableQuantity_FloatingPoint(t *testing.T) {
	transactions := []*model.Transaction{
		{
			AssetID:  1,
			Type:     model.TransactionBuy,
			Quantity: 0.3,
		},
		{
			AssetID:  1,
			Type:     model.TransactionBuy,
			Quantity: 0.2,
		},
	}

	quantity := calculateAvailableQuantity(transactions, 1)
	if quantity < 0.5-1e-9 || quantity > 0.5+1e-9 {
		t.Fatalf("got quantity %v, want approximately 0.5", quantity)
	}
}

func TestTransactionService_Create_SellAfterPositionClosed(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()

	asset := &model.Asset{
		Symbol: "BTC",
		Name:   "Bitcoin",
	}

	if err := assetRepo.Create(context.Background(), asset); err != nil {
		t.Fatal(err)
	}

	service := NewTransactionService(repo, assetRepo)

	buy := &model.Transaction{
		UserID:   1,
		AssetID:  asset.ID,
		Type:     model.TransactionBuy,
		Quantity: 1,
		Price:    100,
	}

	if err := service.Create(context.Background(), buy); err != nil {
		t.Fatal(err)
	}

	sell1 := &model.Transaction{
		UserID:   1,
		AssetID:  asset.ID,
		Type:     model.TransactionSell,
		Quantity: 1,
		Price:    120,
	}

	if err := service.Create(context.Background(), sell1); err != nil {
		t.Fatal(err)
	}

	sell2 := &model.Transaction{
		UserID:   1,
		AssetID:  asset.ID,
		Type:     model.TransactionSell,
		Quantity: 0.1,
		Price:    120,
	}

	err := service.Create(context.Background(), sell2)
	if err == nil {
		t.Fatal("expected insufficient quantity, got nil")
	}
	if err.Error() != "insufficient quantity" {
		t.Fatalf("got error %q, want %q", err.Error(), "insufficient quantity")
	}
}

func TestTransactionService_Create_SellOtherUserQuantity(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()

	asset := &model.Asset{
		Symbol: "BTC",
		Name:   "Bitcoin",
	}

	if err := assetRepo.Create(context.Background(), asset); err != nil {
		t.Fatal(err)
	}

	service := NewTransactionService(repo, assetRepo)

	buy := &model.Transaction{
		UserID:   1,
		AssetID:  asset.ID,
		Type:     model.TransactionBuy,
		Quantity: 1,
		Price:    100,
	}

	if err := service.Create(context.Background(), buy); err != nil {
		t.Fatal(err)
	}

	sell := &model.Transaction{
		UserID:   2,
		AssetID:  asset.ID,
		Type:     model.TransactionSell,
		Quantity: 1,
		Price:    120,
	}

	err := service.Create(context.Background(), sell)

	if err == nil {
		t.Fatal("expected insufficient quantity, got nil")
	}
	if err.Error() != "insufficient quantity" {
		t.Fatalf("got error %q, want %q", err.Error(), "insufficient quantity")
	}
}

func TestTransactionService_Create_Commission(t *testing.T) {
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
		Commission: 2.5,
	}

	err := service.Create(context.Background(), transaction)
	if err != nil {
		t.Fatalf("failed to create transaction: %v", err)
	}

	got, err := repo.FindByID(context.Background(), transaction.ID)
	if err != nil {
		t.Errorf("got error %v, want %v", err, repository.ErrTransactionNotFound)
	}
	if got.Commission != transaction.Commission {
		t.Errorf("got commission %f, want %f", got.Commission, transaction.Commission)
	}
}

func TestTransactionService_Create_SellAfterBuyAgain(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()

	asset := &model.Asset{
		Symbol: "BTC",
		Name:   "Bitcoin",
	}

	if err := assetRepo.Create(context.Background(), asset); err != nil {
		t.Fatal(err)
	}

	service := NewTransactionService(repo, assetRepo)

	buy1 := &model.Transaction{
		UserID:   1,
		AssetID:  asset.ID,
		Type:     model.TransactionBuy,
		Quantity: 1,
		Price:    100,
	}

	if err := service.Create(context.Background(), buy1); err != nil {
		t.Fatal(err)
	}

	sell1 := &model.Transaction{
		UserID:   1,
		AssetID:  asset.ID,
		Type:     model.TransactionSell,
		Quantity: 0.6,
		Price:    120,
	}

	if err := service.Create(context.Background(), sell1); err != nil {
		t.Fatal(err)
	}

	buy2 := &model.Transaction{
		UserID:   1,
		AssetID:  asset.ID,
		Type:     model.TransactionBuy,
		Quantity: 0.5,
		Price:    100,
	}

	if err := service.Create(context.Background(), buy2); err != nil {
		t.Fatal(err)
	}

	sell2 := &model.Transaction{
		UserID:   1,
		AssetID:  asset.ID,
		Type:     model.TransactionSell,
		Quantity: 0.9,
		Price:    120,
	}

	if err := service.Create(context.Background(), sell2); err != nil {
		t.Fatal(err)
	}

	transactions, err := repo.FindByUserID(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}

	if len(transactions) != 4 {
		t.Fatalf("expected 4 transaction, got %d", len(transactions))
	}
}

func TestTransactionService_Create_AssetNotFound(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()
	service := NewTransactionService(repo, assetRepo)

	transaction := &model.Transaction{
		UserID:   1,
		AssetID:  999,
		Type:     model.TransactionBuy,
		Quantity: 1,
		Price:    100,
	}

	err := service.Create(context.Background(), transaction)

	if err == nil {
		t.Fatal("expected asset not found error, got nil")
	}

	transactions, err := repo.FindByUserID(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}

	if len(transactions) != 0 {
		t.Fatalf("expected 0 transactions, got %d", len(transactions))
	}
}

func TestTransactionService_Create_SellCommission(t *testing.T) {
	repo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()

	asset := &model.Asset{
		Symbol: "BTC",
		Name:   "Bitcoin",
	}

	if err := assetRepo.Create(context.Background(), asset); err != nil {
		t.Fatal(err)
	}

	service := NewTransactionService(repo, assetRepo)

	buy1 := &model.Transaction{
		UserID:   1,
		AssetID:  asset.ID,
		Type:     model.TransactionBuy,
		Quantity: 1,
		Price:    100,
	}

	if err := service.Create(context.Background(), buy1); err != nil {
		t.Fatal(err)
	}

	sell1 := &model.Transaction{
		UserID:     1,
		AssetID:    asset.ID,
		Type:       model.TransactionSell,
		Quantity:   0.6,
		Price:      120,
		Commission: 2.5,
	}

	if err := service.Create(context.Background(), sell1); err != nil {
		t.Fatal(err)
	}

	created, err := repo.FindByID(context.Background(), sell1.ID)
	if err != nil {
		t.Fatal(err)
	}

	if created.Commission != 2.5 {
		t.Fatalf("expected commision 2.5, got %f", created.Commission)
	}
}
