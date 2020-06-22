package domain

import "time"

// 取引
type Trades struct {
	ID           int
	TradeID      int `toml:"column:trade_id"`
	Units        float64
	Price        float64
	Instrument   string
	State        string
	InitialUnits float64
	CurrentUnits float64
	RealizedPl   float64
	UnrealizedPl float64
	MarginUsed   float64
	OpenTime     time.Time
	CloseTime    time.Time
}
