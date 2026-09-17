package service

import (
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

func TestCalculateProfitLoss(t *testing.T) {
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
			got := CalculateProfitLoss(tt.quantity, tt.averagePrice, tt.currentPrice)
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

func TestCalculateAveragePrice(t *testing.T) {
	tests := []struct {
		name      string
		quantity1 float64
		price1    float64
		quantity2 float64
		price2    float64
		expected  float64
	}{
		{name: "equal_quantity", quantity1: 10, price1: 180, quantity2: 10, price2: 200, expected: 190},
		{name: "different_quantity", quantity1: 10, price1: 180, quantity2: 5, price2: 200, expected: 186.666667},
		{name: "single_position", quantity1: 10, price1: 180, quantity2: 0, price2: 0, expected: 180},
		{name: "zero_quantity", quantity1: 0, price1: 180, quantity2: 0, price2: 200, expected: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateAveragePrice(tt.quantity1, tt.price1, tt.quantity2, tt.price2)
			if math.Abs(got-tt.expected) >= 0.000001 {
				t.Errorf("got %.6f, want %.6f", got, tt.expected)
			}
		})
	}
}
