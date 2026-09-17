package handler

import (
	"PFnPTA/internal/model"
	"PFnPTA/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type TransactionHandler struct {
	transactionService *service.TransactionService
}

type CreateTransactionRequest struct {
	UserID     int64                 `json:"userId"`
	AssetID    int64                 `json:"assetId"`
	Type       model.TransactionType `json:"type"`
	Quantity   float64               `json:"quantity"`
	Price      float64               `json:"price"`
	Commission float64               `json:"commission"`
}

func NewTransactionHandler(transactionService *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService}
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateTransactionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	transaction := &model.Transaction{
		UserID:     req.UserID,
		AssetID:    req.AssetID,
		Type:       req.Type,
		Quantity:   req.Quantity,
		Price:      req.Price,
		Commission: req.Commission,
	}

	if err := h.transactionService.Create(r.Context(), transaction); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(transaction); err != nil {
		return
	}
}

func (h *TransactionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/transactions/")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid transaction id", http.StatusBadRequest)
		return
	}

	transaction, err := h.transactionService.FindByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(transaction); err != nil {
		return
	}
}

func (h *TransactionHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	userIDStr := query.Get("userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)

	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	transaction, err := h.transactionService.FindByUserID(
		r.Context(),
		userID,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(transaction); err != nil {
		return
	}
}

func (h *TransactionHandler) HandleCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetByUserID(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
