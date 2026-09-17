package model

import "time"

type TransactionType string

const (
	TransactionBuy  TransactionType = "BUY"
	TransactionSell TransactionType = "SELL"
)

type Transaction struct {
	ID         int64
	UserID     int64
	AssetID    int64
	Type       TransactionType
	Quantity   float64
	Price      float64
	Commission float64
	CreatedAt  time.Time
}
