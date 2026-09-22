package main

import (
	"PFnPTA/internal/client"
	"PFnPTA/internal/config"
	"PFnPTA/internal/handler"
	"PFnPTA/internal/middleware"
	"PFnPTA/internal/repository"
	"PFnPTA/internal/service"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func buildMux(userHandler *handler.UserHandler, assetHandler *handler.AssetHandler, authMiddleware *middleware.AuthMiddleware, transactionHandler *handler.TransactionHandler, portfolioHandler *handler.PortfolioHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handler.HealthHandler)
	mux.HandleFunc("/register", userHandler.Register)
	mux.HandleFunc("/login", userHandler.Login)
	mux.HandleFunc("/assets", assetHandler.HandelCollection)
	mux.HandleFunc("/assets/", assetHandler.HandelByID)
	mux.Handle("/transactions", authMiddleware.RequireAuth(http.HandlerFunc(transactionHandler.HandleCollection)))
	mux.Handle("/transactions/", authMiddleware.RequireAuth(http.HandlerFunc(transactionHandler.GetByID)))
	mux.Handle("/me", authMiddleware.RequireAuth(http.HandlerFunc(userHandler.Me)))
	mux.Handle("/portfolio", authMiddleware.RequireAuth(http.HandlerFunc(portfolioHandler.GetPortfolio)))

	return mux
}

func main() {
	log.Println("portfolio tracker starting...")

	cfg := config.Load()

	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	marketData := client.NewBinanceMarketDataProvider(httpClient)

	ctx := context.Background()
	userRepository, err := repository.NewPostgresRepository(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer userRepository.Close(ctx)
	userService := service.NewUserService(userRepository)
	jwtService := service.NewJWTService(cfg.JWTSecret)
	authMiddleware := middleware.NewAuthMiddleware(jwtService)
	userHandler := handler.NewUserHandler(userService, jwtService)

	assetRepository := repository.NewPostgresAssetRepository(userRepository)
	assetService := service.NewAssetService(assetRepository)
	assetHandler := handler.NewAssetHandler(assetService)

	transactionRepository := repository.NewPostgresTransactionRepository(userRepository)
	transactionService := service.NewTransactionService(transactionRepository, assetRepository)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	portfolioService := service.NewPortfolioService(transactionRepository, assetRepository, marketData)
	portfolioHandler := handler.NewPortfolioHandler(portfolioService)

	mux := buildMux(userHandler, assetHandler, authMiddleware, transactionHandler, portfolioHandler)

	server := http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(
		stop,
		os.Interrupt,
	)

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("An error occurred while starting the server: %v", err)
		} else {
			log.Printf("Server was successfully stopped")
		}
	}()

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Failed to shut down the server correctly: %v", err)
	}
}
