package domain

import "github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"

// 足データ（Bid/Ask）
type BidAskCandles struct {
	ID             int
	InstrumentName string
	Granularity    enum.Granularity
	Bid            BidRate `gorm:"embedded"`
	Ask            AskRate `gorm:"embedded"`
	Candles        Candles `gorm:"embedded"`
	Line           enum.Line
	Trend          enum.Trend
}
