package model

import "time"

type Asset struct {
	ID        int64
	Symbol    string
	Name      string
	CreatedAt time.Time
}
