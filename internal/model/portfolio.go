package model

type PortfolioPosition struct {
	AssetID            int64
	Quantity           float64
	AveragePrice       float64
	CurrentPrice       float64
	PositionValue      float64
	ProfitLoss         float64
	ProfitLossPercent  float64
	RealizedProfitLoss float64
}
