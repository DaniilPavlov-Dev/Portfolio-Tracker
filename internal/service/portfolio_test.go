package service

import (
	"PFnPTA/internal/model"
	"math"
	"testing"
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

	position := BuildPosition(transactions)

	if position.Quantity != 0 {
		t.Errorf("expected quantity 0, got %v", position.Quantity)
	}

	if math.Abs(position.AveragePrice-134.666667) >= 0.000001 {
		t.Errorf("expected average price ~134.666667, got %v", position.AveragePrice)
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
