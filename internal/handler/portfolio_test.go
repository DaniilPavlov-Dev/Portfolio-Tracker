package handler

import (
	"PFnPTA/internal/middleware"
	"PFnPTA/internal/model"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockPortfolioService struct {
	positions []*model.PortfolioPosition
	err       error
	userID    int64
}

func (m *mockPortfolioService) GetPortfolio(ctx context.Context, userID int64) ([]*model.PortfolioPosition, error) {
	m.userID = userID
	return m.positions, m.err
}

func TestPortfolioHandler_GetPortfolio(t *testing.T) {
	mockService := &mockPortfolioService{
		positions: []*model.PortfolioPosition{
			{
				AssetID:       1,
				Quantity:      2,
				AveragePrice:  100,
				CurrentPrice:  120,
				PositionValue: 240,
				ProfitLoss:    40,
			},
		},
	}

	handler := NewPortfolioHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/portfolio", nil)

	userID := int64(42)
	req = req.WithContext(middleware.ContextWithUserID(req.Context(), userID))

	rec := httptest.NewRecorder()

	handler.GetPortfolio(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", rec.Code, http.StatusOK)
	}

	if mockService.userID != userID {
		t.Fatalf("expected userID %d, got %d", userID, mockService.userID)
	}

	var positions []*model.PortfolioPosition

	if err := json.NewDecoder(rec.Body).Decode(&positions); err != nil {
		t.Fatalf("expected to decode response: %v", err)
	}

	if len(positions) != 1 {
		t.Fatalf("expected 1 position, got %d", len(positions))
	}

	if positions[0].AssetID != 1 {
		t.Fatalf("expected assetID 1, got %d", positions[0].AssetID)
	}

	if positions[0].Quantity != 2 {
		t.Fatalf("expected quantity 2, got %f", positions[0].Quantity)
	}
}

func TestPortfolioHandler_GetPortfolio_Unauthorized(t *testing.T) {
	mockService := &mockPortfolioService{}

	handler := NewPortfolioHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/portfolio", nil)
	rec := httptest.NewRecorder()

	handler.GetPortfolio(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusUnauthorized,
			rec.Code,
		)
	}
}
