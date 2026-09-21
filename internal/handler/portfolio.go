package handler

import (
	"PFnPTA/internal/middleware"
	"PFnPTA/internal/model"
	"context"
	"encoding/json"
	"net/http"
)

type PortfolioHandler struct {
	portfolioService PortfolioService
}

type PortfolioService interface {
	GetPortfolio(ctx context.Context, userID int64) ([]*model.PortfolioPosition, error)
}

func NewPortfolioHandler(portfolioService PortfolioService) *PortfolioHandler {
	return &PortfolioHandler{portfolioService: portfolioService}
}

func (h *PortfolioHandler) GetPortfolio(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unathorized", http.StatusUnauthorized)
		return
	}

	positions, err := h.portfolioService.GetPortfolio(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to get portfolio", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(positions); err != nil {
		http.Error(w, "falied to encode positions", http.StatusInternalServerError)
	}
}
