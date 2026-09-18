package service

import (
	"PFnPTA/internal/model"
)

type portfolioState struct {
	quantity           float64
	averagePrice       float64
	realizedProfitLoss float64
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

func BuildPosition(transactions []*model.Transaction) *model.PortfolioPosition {
	state := calculatePortfolioState(transactions)

	return &model.PortfolioPosition{Quantity: state.quantity, AveragePrice: state.averagePrice}
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
