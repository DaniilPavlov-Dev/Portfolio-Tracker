package client

import (
	"PFnPTA/internal/service"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type binancePriceResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

type StaticMarketDataProvider struct {
	prices map[string]float64
}

var _ service.MarketDataProvider = (*BinanceMarketDataProvider)(nil)

var _ service.MarketDataProvider = (*StaticMarketDataProvider)(nil)

type MarketDataProvider interface {
	GetPrice(ctx context.Context, symbol string) (float64, error)
}

type BinanceMarketDataProvider struct {
	client  *http.Client
	baseURL string
}

func NewBinanceMarketDataProvider(client *http.Client) *BinanceMarketDataProvider {
	return &BinanceMarketDataProvider{client: client, baseURL: "https://api.binance.com"}
}

func NewStaticMarketDataProvider() *StaticMarketDataProvider {
	return &StaticMarketDataProvider{
		prices: map[string]float64{
			"BTC": 100000,
			"ETH": 4000,
		},
	}
}

func (p *StaticMarketDataProvider) GetPrice(ctx context.Context, symbol string) (float64, error) {
	price, ok := p.prices[symbol]
	if !ok {
		return 0, nil
	}

	return price, nil
}

func (p *BinanceMarketDataProvider) GetPrice(ctx context.Context, symbol string) (float64, error) {
	symbol = symbol + "USDT"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.baseURL+"/api/v3/ticker/price?symbol="+symbol, nil)
	if err != nil {
		return 0, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Binance returned status: %d", resp.StatusCode)
	}

	var data binancePriceResponse

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	price, err := strconv.ParseFloat(data.Price, 64)
	if err != nil {
		return 0, err
	}

	return price, nil
}
