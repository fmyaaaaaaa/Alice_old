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
	SwingID        int
}

func (b BidAskCandles) GetOpenMid() float64 {
	return (b.Bid.Open + b.Ask.Open) / 2
}

func (b BidAskCandles) GetCloseMid() float64 {
	return (b.Bid.Close + b.Ask.Close) / 2
}

func (b BidAskCandles) GetHighMid() float64 {
	return (b.Bid.High + b.Ask.High) / 2
}

func (b BidAskCandles) GetLowMid() float64 {
	return (b.Bid.Low + b.Ask.Low) / 2
}

func (b BidAskCandles) GetAveMid() float64 {
	return (b.GetOpenMid() + b.GetCloseMid()) / 2
}
