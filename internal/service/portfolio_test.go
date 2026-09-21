package service

import (
	"PFnPTA/internal/model"
	"PFnPTA/internal/repository"
	"context"
	"errors"
	"math"
	"testing"
	"time"
)

func TestCalculatePositionValue(t *testing.T) {
	tests := []struct {
		name         string
		quantity     float64
		currentPrice float64
		expected     float64
	}{
		{name: "regular", quantity: 10, currentPrice: 180, expected: 1800},
		{name: "bitcoin", quantity: 0.5, currentPrice: 60000, expected: 30000},
		{name: "zero", quantity: 0, currentPrice: 100, expected: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculatePositionValue(tt.quantity, tt.currentPrice)
			if got != tt.expected {
				t.Errorf("got %.2f, want %.2f", got, tt.expected)
			}
		})
	}
}

func TestCalculateUnrealizedProfitLoss(t *testing.T) {
	tests := []struct {
		name         string
		quantity     float64
		averagePrice float64
		currentPrice float64
		expected     float64
	}{
		{name: "regular", quantity: 10, averagePrice: 180, currentPrice: 195, expected: 150},
		{name: "loss", quantity: 10, averagePrice: 180, currentPrice: 170, expected: -100},
		{name: "zero", quantity: 0, averagePrice: 180, currentPrice: 195, expected: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateUnrealizedProfitLoss(tt.quantity, tt.averagePrice, tt.currentPrice)
			if got != tt.expected {
				t.Errorf("got %.2f, want %.2f", got, tt.expected)
			}
		})
	}
}

func TestCalculateProfitLossPercent(t *testing.T) {
	tests := []struct {
		name         string
		quantity     float64
		averagePrice float64
		currentPrice float64
		expected     float64
	}{
		{name: "profit", quantity: 10, averagePrice: 180, currentPrice: 195, expected: 8.333333},
		{name: "loss", quantity: 10, averagePrice: 180, currentPrice: 170, expected: -5.555555},
		{name: "same_price", quantity: 10, averagePrice: 180, currentPrice: 180, expected: 0},
		{name: "zero_average", quantity: 10, averagePrice: 0, currentPrice: 195, expected: 0},
		{name: "zero_quantity", quantity: 0, averagePrice: 180, currentPrice: 195, expected: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateProfitLossPercent(tt.quantity, tt.averagePrice, tt.currentPrice)
			if math.Abs(got-tt.expected) >= 0.000001 {
				t.Errorf("got %.6f, want %.6f", got, tt.expected)
			}
		})
	}
}

func TestBuildPosition(t *testing.T) {
	transactions := []*model.Transaction{
		{
			Type:       model.TransactionBuy,
			Quantity:   10,
			Price:      100,
			Commission: 20,
		},
		{
			Type:     model.TransactionBuy,
			Quantity: 5,
			Price:    200,
		},
		{
			Type:     model.TransactionSell,
			Quantity: 3,
			Price:    150,
		},
		{
			Type:     model.TransactionSell,
			Quantity: 12,
			Price:    180,
		},
	}

	position := BuildPosition(10, transactions, 150)

	if position != nil {
		t.Errorf("expected nil position, got %+v", position)
	}
}

func TestBuildPosition_Empty(t *testing.T) {
	position := BuildPosition(10, nil, 150)

	if position != nil {
		t.Errorf("expected nil position, got %+v", position)
	}
}

func TestCalculateRealizedProfitLoss(t *testing.T) {
	transactions := []*model.Transaction{
		{
			Type:       model.TransactionBuy,
			Quantity:   10,
			Price:      100,
			Commission: 20,
		},
		{
			Type:     model.TransactionBuy,
			Quantity: 5,
			Price:    200,
		},
		{
			Type:       model.TransactionSell,
			Quantity:   3,
			Price:      150,
			Commission: 10,
		},
		{
			Type:     model.TransactionSell,
			Quantity: 12,
			Price:    180,
		},
	}

	got := CalculateRealizedProfitLoss(transactions)

	if math.Abs(got-580) >= 0.000001 {
		t.Errorf("expected realized P/L 580, got %v", got)
	}
}

func TestCalculateBuyAveragePrice(t *testing.T) {
	got := CalculateBuyAveragePrice(0, 0, 10, 100, 20)

	if math.Abs(got-102) >= 0.000001 {
		t.Errorf("expected average price 102, got %v", got)
	}
}

func TestCalculateRealizedProfitLoss_Loss(t *testing.T) {
	transactions := []*model.Transaction{
		{
			Type:     model.TransactionBuy,
			Quantity: 10,
			Price:    100,
		},
		{
			Type:     model.TransactionSell,
			Quantity: 4,
			Price:    80,
		},
	}

	got := CalculateRealizedProfitLoss(transactions)

	if math.Abs(got-(-80)) >= 0.000001 {
		t.Errorf("expected realized P/L -80, got %v", got)
	}
}

func TestCalculateRealizedProfitLoss_PartialSellAndBuy(t *testing.T) {
	transactions := []*model.Transaction{
		{
			Type:     model.TransactionBuy,
			Quantity: 10,
			Price:    100,
		},
		{
			Type:     model.TransactionSell,
			Quantity: 4,
			Price:    150,
		},
		{
			Type:     model.TransactionBuy,
			Quantity: 6,
			Price:    200,
		},
	}

	got := CalculateRealizedProfitLoss(transactions)
	if math.Abs(got-200) >= 0.000001 {
		t.Errorf("expected realized P/L 200, got %v", got)
	}
}

func TestCalculatePosition(t *testing.T) {
	tests := []struct {
		name                  string
		quantity              float64
		averagePrice          float64
		currentPrice          float64
		expectedValue         float64
		expectedProfit        float64
		expectedProfitPercent float64
	}{
		{
			name:                  "profit",
			quantity:              10,
			averagePrice:          100,
			currentPrice:          120,
			expectedValue:         1200,
			expectedProfit:        200,
			expectedProfitPercent: 20,
		},
		{
			name:                  "loss",
			quantity:              10,
			averagePrice:          100,
			currentPrice:          80,
			expectedValue:         800,
			expectedProfit:        -200,
			expectedProfitPercent: -20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			position := CalculatePosition(
				tt.quantity,
				tt.averagePrice,
				tt.currentPrice,
			)

			if math.Abs(position.PositionValue-tt.expectedValue) >= 0.000001 {
				t.Errorf("expected position value %.2f, got %.2f", tt.expectedValue, position.PositionValue)
			}

			if math.Abs(position.ProfitLoss-tt.expectedProfit) >= 0.000001 {
				t.Errorf("expected profit/loss %.2f, got %.2f", tt.expectedProfit, position.ProfitLoss)
			}

			if math.Abs(position.ProfitLossPercent-tt.expectedProfitPercent) >= 0.000001 {
				t.Errorf("expected profit/loss percent %.2f, got %.2f", tt.expectedProfitPercent, position.ProfitLossPercent)
			}
		})
	}
}

type mockMarketDataProvider struct {
	prices map[string]float64
}

type errorMarketDataProvider struct {
}

func (m *errorMarketDataProvider) GetPrice(ctx context.Context, symbol string) (float64, error) {
	return 0, errors.New("market data error")
}

func (m *mockMarketDataProvider) GetPrice(ctx context.Context, symbol string) (float64, error) {
	return m.prices[symbol], nil
}
func TestPortfolioService_GetTransactions(t *testing.T) {
	transactionRepo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()
	marketData := &mockMarketDataProvider{}

	portfolioService := NewPortfolioService(transactionRepo, assetRepo, marketData)

	transaction := &model.Transaction{
		UserID:     1,
		AssetID:    10,
		Type:       model.TransactionBuy,
		Quantity:   2,
		Price:      100,
		Commission: 1,
	}

	err := transactionRepo.Create(context.Background(), transaction)
	if err != nil {
		t.Fatalf("failed to create transaction: %v", err)
	}

	transactions, err := portfolioService.GetTransactions(context.Background(), 1)
	if err != nil {
		t.Fatalf("failed to get transaction: %v", err)
	}

	if len(transactions) != 1 {
		t.Fatalf("expected 1 transsation, got %d", len(transactions))
	}

	if transactions[0].ID != transaction.ID {
		t.Errorf("expected transactiion ID %d, got %d", transaction.ID, transactions[0].ID)
	}
}

func TestGroupTransactionsByAsset(t *testing.T) {
	transactions := []*model.Transaction{
		{
			AssetID: 10,
			Type:    model.TransactionBuy,
		},
		{
			AssetID: 10,
			Type:    model.TransactionSell,
		},
		{
			AssetID: 20,
			Type:    model.TransactionBuy,
		},
	}

	grouped := groupTransactionsByAsset(transactions)

	if len(grouped[10]) != 2 {
		t.Errorf("expected 2 transactions for asset 10, got %d", len(grouped[10]))
	}

	if len(grouped[20]) != 1 {
		t.Errorf("expected 1 transactions for asset 20, got %d", len(grouped[20]))
	}
}

func TestPortfolioService_GetPortfolio(t *testing.T) {
	transactionRepo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()
	marketData := &mockMarketDataProvider{prices: map[string]float64{"BTC": 120, "ETH": 300}}

	portfolioService := NewPortfolioService(transactionRepo, assetRepo, marketData)

	asset := &model.Asset{
		Symbol: "BTC",
		Name:   "Bitcoin",
	}

	err := assetRepo.Create(context.Background(), asset)
	if err != nil {
		t.Fatal(err)
	}

	eth := &model.Asset{
		Symbol: "ETH",
		Name:   "Ethereum",
	}

	err = assetRepo.Create(context.Background(), eth)
	if err != nil {
		t.Fatal(err)
	}

	transaction := &model.Transaction{
		UserID:     1,
		AssetID:    asset.ID,
		Type:       model.TransactionBuy,
		Quantity:   10,
		Price:      100,
		Commission: 0,
	}

	err = transactionRepo.Create(context.Background(), transaction)
	if err != nil {
		t.Fatal(err)
	}

	ethTransaction := &model.Transaction{
		UserID:     1,
		AssetID:    eth.ID,
		Type:       model.TransactionBuy,
		Quantity:   5,
		Price:      200,
		Commission: 0,
	}

	err = transactionRepo.Create(context.Background(), ethTransaction)
	if err != nil {
		t.Fatal(err)
	}

	positions, err := portfolioService.GetPortfolio(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}

	if len(positions) != 2 {
		t.Errorf("expected 2 position, got %d", len(positions))
	}

	var ethPosition *model.PortfolioPosition

	for _, position := range positions {
		if position.AssetID == eth.ID {
			ethPosition = position
			break
		}
	}

	if ethPosition == nil {
		t.Fatal("expected ETH position")
	}

	var btcPosition *model.PortfolioPosition

	for _, position := range positions {
		if position.AssetID == asset.ID {
			btcPosition = position
			break
		}
	}

	if btcPosition == nil {
		t.Fatal("expected BTC position")
	}

	position := btcPosition
	if position.Quantity != 10 {
		t.Errorf("expected quantity 10, got %f", position.Quantity)
	}

	if position.AveragePrice != 100 {
		t.Errorf("expected average price 100, got %f", position.AveragePrice)
	}

	if position.CurrentPrice != 120 {
		t.Errorf("expected current price 120, got %f", position.CurrentPrice)
	}

	if position.PositionValue != 1200 {
		t.Errorf("expected position value 1200, got %f", position.PositionValue)
	}

	if position.ProfitLoss != 200 {
		t.Errorf("expected profit/loss 200, got %f", position.ProfitLoss)
	}

	if position.ProfitLossPercent != 20 {
		t.Errorf("expected profit/loss percent 20, got %f", position.ProfitLossPercent)
	}

	if ethPosition.Quantity != 5 {
		t.Errorf("expected ETH quantity 5, got %f", ethPosition.Quantity)
	}

	if ethPosition.AveragePrice != 200 {
		t.Errorf("expected ETH average price 200, got %f", ethPosition.AveragePrice)
	}

	if ethPosition.CurrentPrice != 300 {
		t.Errorf("expected current price 300, got %f", ethPosition.CurrentPrice)
	}

	if ethPosition.PositionValue != 1500 {
		t.Errorf("expected ETH position value 1500, got %f", ethPosition.PositionValue)
	}

	if ethPosition.ProfitLoss != 500 {
		t.Errorf("expected ETH profit/loss 500, got %f", ethPosition.ProfitLoss)
	}

	if ethPosition.ProfitLossPercent != 50 {
		t.Errorf("expected ETH profit/loss percent 50, got %f", ethPosition.ProfitLossPercent)
	}

	ethSellTransaction := &model.Transaction{
		UserID:     1,
		AssetID:    eth.ID,
		Type:       model.TransactionSell,
		Quantity:   2,
		Price:      250,
		Commission: 0,
	}

	err = transactionRepo.Create(context.Background(), ethSellTransaction)
	if err != nil {
		t.Fatal(err)
	}

	sellTransaction := &model.Transaction{
		UserID:     1,
		AssetID:    asset.ID,
		Type:       model.TransactionSell,
		Quantity:   10,
		Price:      110,
		Commission: 0,
	}

	err = transactionRepo.Create(context.Background(), sellTransaction)
	if err != nil {
		t.Fatal(err)
	}

	positions, err = portfolioService.GetPortfolio(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}

	if len(positions) != 1 {
		t.Fatalf("expected 1 position after selling BTC, got %d", len(positions))
	}

	if positions[0].AssetID != eth.ID {
		t.Errorf(
			"expected ETH position after selling BTC, got asset ID %d",
			positions[0].AssetID,
		)
	}

	if positions[0].RealizedProfitLoss != 100 {
		t.Errorf(
			"expected realized profit/loss 100, got %f",
			positions[0].RealizedProfitLoss,
		)
	}
}

func TestPortfolioService_GetPortfolio_MarketDataError(t *testing.T) {
	transactionRepo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()
	marketData := &errorMarketDataProvider{}

	portfolioService := NewPortfolioService(transactionRepo, assetRepo, marketData)

	asset := &model.Asset{
		Symbol: "BTC",
		Name:   "Bitcoin",
	}

	err := assetRepo.Create(context.Background(), asset)
	if err != nil {
		t.Fatal(err)
	}

	transaction := &model.Transaction{
		UserID:     1,
		AssetID:    asset.ID,
		Type:       model.TransactionBuy,
		Quantity:   10,
		Price:      100,
		Commission: 0,
	}

	err = transactionRepo.Create(context.Background(), transaction)
	if err != nil {
		t.Fatal(err)
	}

	_, err = portfolioService.GetPortfolio(
		context.Background(),
		1,
	)
	if err == nil {
		t.Fatal("expected market data error")
	}
}

func TestPortfolioService_GetPortfolio_Empty(t *testing.T) {
	transactionRepo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()
	marketData := &mockMarketDataProvider{
		prices: map[string]float64{},
	}

	portfolioService := NewPortfolioService(
		transactionRepo,
		assetRepo,
		marketData,
	)

	positions, err := portfolioService.GetPortfolio(
		context.Background(),
		1,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(positions) != 0 {
		t.Fatalf("expected 0 positions, got %d", len(positions))
	}
}

func TestPortfolioService_GetPortfolio_AssetNotFound(t *testing.T) {
	transactionRepo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()
	marketData := &mockMarketDataProvider{
		prices: map[string]float64{},
	}

	portfolioService := NewPortfolioService(
		transactionRepo,
		assetRepo,
		marketData,
	)

	transaction := &model.Transaction{
		UserID:     1,
		AssetID:    999,
		Type:       model.TransactionBuy,
		Quantity:   10,
		Price:      100,
		Commission: 0,
	}

	err := transactionRepo.Create(
		context.Background(),
		transaction,
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = portfolioService.GetPortfolio(
		context.Background(),
		1,
	)
	if err == nil {
		t.Fatal("expected asset not found error")
	}
}

func TestBuildPosition_RealizedProfitLoss(t *testing.T) {
	transactions := []*model.Transaction{
		{
			AssetID:    10,
			Type:       model.TransactionBuy,
			Quantity:   5,
			Price:      200,
			Commission: 0,
		},
		{
			AssetID:    10,
			Type:       model.TransactionSell,
			Quantity:   2,
			Price:      250,
			Commission: 0,
		},
	}

	position := BuildPosition(
		10,
		transactions,
		300,
	)

	if position == nil {
		t.Fatal("expected position, got nil")
	}

	if position.RealizedProfitLoss != 100 {
		t.Errorf(
			"expected realized profit/loss 100, got %f",
			position.RealizedProfitLoss,
		)
	}
}

type cancelAwareMarketDataProvider struct {
	canceled chan struct{}
}

func (m *cancelAwareMarketDataProvider) GetPrice(
	ctx context.Context,
	symbol string,
) (float64, error) {
	if symbol == "BTC" {
		return 0, errors.New("market data error")
	}

	<-ctx.Done()

	close(m.canceled)
	return 0, ctx.Err()
}

func TestPortfolioService_GetPortfolio_CancelsOtherRequests(t *testing.T) {
	transactionRepo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()

	canceled := make(chan struct{})

	marketData := &cancelAwareMarketDataProvider{canceled: canceled}

	portfolioService := NewPortfolioService(transactionRepo, assetRepo, marketData)

	btc := &model.Asset{
		Symbol: "BTC",
		Name:   "Bitcoin",
	}

	eth := &model.Asset{
		Symbol: "ETH",
		Name:   "Ethereum",
	}

	if err := assetRepo.Create(context.Background(), btc); err != nil {
		t.Fatal(err)
	}
	if err := assetRepo.Create(context.Background(), eth); err != nil {
		t.Fatal(err)
	}

	transactions := []*model.Transaction{
		{
			UserID:   1,
			AssetID:  btc.ID,
			Type:     model.TransactionBuy,
			Quantity: 1,
			Price:    100,
		},
		{
			UserID:   1,
			AssetID:  eth.ID,
			Type:     model.TransactionBuy,
			Quantity: 1,
			Price:    100,
		},
	}

	for _, transaction := range transactions {
		if err := transactionRepo.Create(context.Background(), transaction); err != nil {
			t.Fatal(err)
		}
	}

	_, err := portfolioService.GetPortfolio(context.Background(), 1)

	if err == nil {
		t.Fatal("expected market data error")
	}

	select {
	case <-canceled:
	case <-time.After(time.Second):
		t.Fatal("expected cancellation")
	}
}

type slowMarketDataProvider struct {
	delay time.Duration
}

func (m *slowMarketDataProvider) GetPrice(ctx context.Context, symbol string) (float64, error) {
	select {
	case <-time.After(m.delay):
		return 100, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

func TestPortfolioService_GetPortfolio_ConcurrentMarketData(t *testing.T) {
	transactionRepo := repository.NewMemoryTransactionRepository()
	assetRepo := repository.NewMemoryAssetRepository()

	marketData := &slowMarketDataProvider{
		delay: 500 * time.Millisecond,
	}

	portfolioService := NewPortfolioService(
		transactionRepo,
		assetRepo,
		marketData,
	)

	btc := &model.Asset{
		Symbol: "BTC",
		Name:   "Bitcoin",
	}

	eth := &model.Asset{
		Symbol: "ETH",
		Name:   "Ethereum",
	}

	if err := assetRepo.Create(context.Background(), btc); err != nil {
		t.Fatal(err)
	}

	if err := assetRepo.Create(context.Background(), eth); err != nil {
		t.Fatal(err)
	}

	transactions := []*model.Transaction{
		{
			UserID:   1,
			AssetID:  btc.ID,
			Type:     model.TransactionBuy,
			Quantity: 1,
			Price:    100,
		},
		{
			UserID:   1,
			AssetID:  eth.ID,
			Type:     model.TransactionBuy,
			Quantity: 1,
			Price:    100,
		},
	}

	for _, transaction := range transactions {
		if err := transactionRepo.Create(context.Background(), transaction); err != nil {
			t.Fatal(err)
		}
	}

	start := time.Now()

	_, err := portfolioService.GetPortfolio(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	elapsed := time.Since(start)

	if elapsed >= time.Second {
		t.Fatalf("expected concurrent executuion, took %v", elapsed)
	}
}
