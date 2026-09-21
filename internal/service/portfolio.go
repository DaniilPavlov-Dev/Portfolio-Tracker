package service

import (
	"PFnPTA/internal/model"
	"PFnPTA/internal/repository"
	"context"
)

type portfolioState struct {
	quantity           float64
	averagePrice       float64
	realizedProfitLoss float64
}

type PortfolioService struct {
	transactionRepo repository.TransactionRepository
	assetRepo       repository.AssetRepository
	marketData      MarketDataProvider
}

type MarketDataProvider interface {
	GetPrice(ctx context.Context, symbol string) (float64, error)
}

func NewPortfolioService(transactionRepo repository.TransactionRepository, assetRepo repository.AssetRepository, marketData MarketDataProvider) *PortfolioService {
	return &PortfolioService{transactionRepo: transactionRepo, assetRepo: assetRepo, marketData: marketData}
}

func CalculatePositionValue(quantity, currentPrice float64) float64 {
	return quantity * currentPrice
}

func CalculateProfitLossPercent(quantity, averagePrice, currentPrice float64) float64 {
	if averagePrice == 0 || quantity == 0 {
		return 0
	}
	return (currentPrice - averagePrice) / averagePrice * 100
}

func BuildPosition(assetID int64, transactions []*model.Transaction, currentPrice float64) *model.PortfolioPosition {
	state := calculatePortfolioState(transactions)

	if state.quantity == 0 {
		return nil
	}

	position := CalculatePosition(state.quantity, state.averagePrice, currentPrice)
	position.AssetID = assetID

	position.RealizedProfitLoss = CalculateRealizedProfitLoss(transactions)

	return position
}

func CalculateRealizedProfitLoss(transactions []*model.Transaction) float64 {
	state := calculatePortfolioState(transactions)

	return state.realizedProfitLoss
}

func CalculateBuyAveragePrice(quantity1, averagePrice1 float64, quantity2, averagePrice2, commission float64) float64 {
	if quantity1+quantity2 == 0 {
		return 0
	}

	totalCost := quantity1*averagePrice1 + quantity2*averagePrice2 + commission
	return totalCost / (quantity1 + quantity2)
}

func calculatePortfolioState(transactions []*model.Transaction) portfolioState {
	var state portfolioState

	for _, transaction := range transactions {
		switch transaction.Type {
		case model.TransactionBuy:
			state.averagePrice = CalculateBuyAveragePrice(
				state.quantity,
				state.averagePrice,
				transaction.Quantity,
				transaction.Price,
				transaction.Commission,
			)

			state.quantity += transaction.Quantity

		case model.TransactionSell:
			state.realizedProfitLoss += (transaction.Price-state.averagePrice)*transaction.Quantity - transaction.Commission

			state.quantity -= transaction.Quantity
		}
	}

	return state
}

func CalculateUnrealizedProfitLoss(quantity, averagePrice, currentPrice float64) float64 {
	return quantity * (currentPrice - averagePrice)
}

func CalculatePosition(quantity, averagePrice, currentPrice float64) *model.PortfolioPosition {
	return &model.PortfolioPosition{
		Quantity:          quantity,
		AveragePrice:      averagePrice,
		CurrentPrice:      currentPrice,
		PositionValue:     CalculatePositionValue(quantity, currentPrice),
		ProfitLoss:        CalculateUnrealizedProfitLoss(quantity, averagePrice, currentPrice),
		ProfitLossPercent: CalculateProfitLossPercent(quantity, averagePrice, currentPrice),
	}
}

func (s *PortfolioService) GetTransactions(ctx context.Context, userID int64) ([]*model.Transaction, error) {
	return s.transactionRepo.FindByUserID(ctx, userID)
}

func groupTransactionsByAsset(transactions []*model.Transaction) map[int64][]*model.Transaction {
	grouped := make(map[int64][]*model.Transaction)

	for _, transaction := range transactions {
		grouped[transaction.AssetID] = append(grouped[transaction.AssetID], transaction)
	}

	return grouped
}

func (s *PortfolioService) GetPortfolio(ctx context.Context, userID int64) ([]*model.PortfolioPosition, error) {
	transactions, err := s.transactionRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	grouped := groupTransactionsByAsset(transactions)

	positions := make([]*model.PortfolioPosition, 0, len(grouped))

	for assetID := range grouped {
		asset, err := s.assetRepo.FindByID(ctx, assetID)
		if err != nil {
			return nil, err
		}

		currentPrice, err := s.marketData.GetPrice(ctx, asset.Symbol)
		if err != nil {
			return nil, err
		}

		position := BuildPosition(assetID, grouped[assetID], currentPrice)
		if position == nil {
			continue
		}

		positions = append(positions, position)
	}

	return positions, nil
}
