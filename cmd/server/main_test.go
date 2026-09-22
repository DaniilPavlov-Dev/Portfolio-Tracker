package main

import (
	"PFnPTA/internal/client"
	"PFnPTA/internal/handler"
	"PFnPTA/internal/middleware"
	"PFnPTA/internal/model"
	"PFnPTA/internal/repository"
	"PFnPTA/internal/service"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBuildMux(t *testing.T) {
	jwtService := service.NewJWTService("test-secret")
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	userRepository := repository.NewMemoryUserRepository()

	user := &model.User{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: "test-hash",
	}

	if err := userRepository.Create(context.Background(), user); err != nil {
		t.Fatal(err)
	}

	if user.ID != 1 {
		t.Fatalf("got user ID %d, want %d", user.ID, 1)
	}

	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService, jwtService)

	assetRepository := repository.NewMemoryAssetRepository()
	assetService := service.NewAssetService(assetRepository)
	assetHandler := handler.NewAssetHandler(assetService)

	asset := &model.Asset{
		Symbol: "BTC",
		Name:   "Bitcoin",
	}

	if err := assetRepository.Create(context.Background(), asset); err != nil {
		t.Fatal(err)
	}

	transactionRepository := repository.NewMemoryTransactionRepository()
	transactionService := service.NewTransactionService(transactionRepository, assetRepository)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	marketData := client.NewStaticMarketDataProvider()

	portfolioService := service.NewPortfolioService(transactionRepository, assetRepository, marketData)
	portfolioHandler := handler.NewPortfolioHandler(portfolioService)

	mux := buildMux(userHandler, assetHandler, authMiddleware, transactionHandler, portfolioHandler)

	req := httptest.NewRequest(
		http.MethodGet,
		"/health",
		nil,
	)

	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"got status %d, want %d",
			rec.Code,
			http.StatusOK,
		)
	}

	req = httptest.NewRequest(
		http.MethodGet,
		"/me",
		nil,
	)

	rec = httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"got status %d, want %d",
			rec.Code,
			http.StatusUnauthorized,
		)
	}

	token, err := jwtService.GenerateToken(1)
	if err != nil {
		t.Fatal(err)
	}

	req = httptest.NewRequest(
		http.MethodGet,
		"/me",
		nil,
	)

	req.Header.Set(
		"Authorization",
		"Bearer "+token,
	)

	rec = httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"got status %d, want %d",
			rec.Code,
			http.StatusOK,
		)
	}

	request := httptest.NewRequest(http.MethodGet, "/portfolio", nil)
	request.Header.Set("Authorization", "Bearer "+token)

	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", response.Code, http.StatusOK)
	}

	request = httptest.NewRequest(http.MethodGet, "/assets", nil)

	response = httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", response.Code, http.StatusOK)
	}

	request = httptest.NewRequest(http.MethodGet, "/transactions", nil)

	response = httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("got status %d, want %d", response.Code, http.StatusUnauthorized)
	}

	request = httptest.NewRequest(http.MethodGet, "/transactions", nil)
	request.Header.Set("Authorization", "Bearer "+token)

	response = httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", response.Code, http.StatusOK)
	}

	transactionBody := `{
    "assetId": 1,
    "type": "BUY",
    "quantity": 0.5,
    "price": 60000,
    "commission": 10
}`

	request = httptest.NewRequest(
		http.MethodPost,
		"/transactions",
		strings.NewReader(transactionBody),
	)
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-Type", "application/json")

	response = httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("got status %d, want %d", response.Code, http.StatusCreated)
	}

	transactions, err := transactionRepository.FindByUserID(
		context.Background(),
		1,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(transactions) != 1 {
		t.Fatalf("got %d transactions, want %d", len(transactions), 1)
	}

	if transactions[0].AssetID != asset.ID {
		t.Fatalf(
			"got asset ID %d, want %d",
			transactions[0].AssetID,
			asset.ID,
		)
	}

	request = httptest.NewRequest(http.MethodGet, "/transactions", nil)
	request.Header.Set("Authorization", "Bearer "+token)

	response = httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", response.Code, http.StatusOK)
	}

	var transactionsResponse []model.Transaction

	if err := json.NewDecoder(response.Body).Decode(&transactionsResponse); err != nil {
		t.Fatal(err)
	}

	if len(transactionsResponse) != 1 {
		t.Fatalf(
			"got %d transactions, want %d",
			len(transactionsResponse),
			1,
		)
	}

	if transactionsResponse[0].AssetID != asset.ID {
		t.Fatalf(
			"got asset ID %d, want %d",
			transactionsResponse[0].AssetID,
			asset.ID,
		)
	}
}
