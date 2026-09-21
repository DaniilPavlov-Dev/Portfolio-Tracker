package client

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStaticMarketDataProvider_GetPrice(t *testing.T) {
	provider := NewStaticMarketDataProvider()

	price, err := provider.GetPrice(context.Background(), "BTC")
	if err != nil {
		t.Fatalf("exepcted no error, got %v", err)
	}

	if price != 100000 {
		t.Fatalf("expected price 100000, got %f", price)
	}
}

func TestStaticMarketDataProvider_GetPrice_UnknownSymbol(t *testing.T) {
	provider := NewStaticMarketDataProvider()

	price, err := provider.GetPrice(context.Background(), "UNKNOWN")
	if err != nil {
		t.Fatalf("exepcted no error, got %v", err)
	}

	if price != 0 {
		t.Fatalf("expected price 0, got %f", price)
	}
}

func TestBinanceMarketDataProvider_GetPrice(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v3/ticker/price" {
			t.Fatalf("unexpected URL path: %s", r.URL.Path)
		}

		if got := r.URL.Query().Get("symbol"); got != "BTCUSDT" {
			t.Fatalf("unexpected symbol: %s", got)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"symbol":"BTCUSDT","price":"123.45"}`)
	}))
	defer server.Close()

	provider := &BinanceMarketDataProvider{
		client:  server.Client(),
		baseURL: server.URL,
	}

	price, err := provider.GetPrice(context.Background(), "BTC")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if price != 123.45 {
		t.Fatalf("unexpected price: %f", price)
	}
}

func TestBinanceMarketDataProvider_GetPrice_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}))
	defer server.Close()

	provider := &BinanceMarketDataProvider{
		client:  server.Client(),
		baseURL: server.URL,
	}

	_, err := provider.GetPrice(context.Background(), "BTCUSDT")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestBinanceMarketDataProvider_GetPrice_InvalidPrice(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"symbol":"BTCUSDT","price":"not-a-number"}`)
	}))
	defer server.Close()

	provider := &BinanceMarketDataProvider{
		client:  server.Client(),
		baseURL: server.URL,
	}

	_, err := provider.GetPrice(context.Background(), "BTCUSDT")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
