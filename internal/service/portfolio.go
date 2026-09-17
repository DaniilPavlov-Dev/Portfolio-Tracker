package service

func CalculatePositionValue(quantity, currentPrice float64) float64 {
	return quantity * currentPrice
}

func CalculateProfitLoss(quantity, averagePrice, currentPrice float64) float64 {
	return quantity * (currentPrice - averagePrice)
}

func CalculateProfitLossPercent(quantity, averagePrice, currentPrice float64) float64 {
	if averagePrice == 0 || quantity == 0 {
		return 0
	}
	return (currentPrice - averagePrice) / averagePrice * 100
}

func CalculateAveragePrice(quantity1, price1, quantity2, price2 float64) float64 {
	if quantity1+quantity2 == 0 {
		return 0
	}
	return (quantity1*price1 + quantity2*price2) / (quantity1 + quantity2)
}
