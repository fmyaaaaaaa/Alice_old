package domain

import "github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"

// キャプテン・アメリカのセットアップステータス
type CaptainAmericaStatus struct {
	ID          int
	Instrument  string
	Granularity enum.Granularity
	Line        enum.Line
	SetupPrice  float64
	SetupStatus bool
	TradeStatus bool
	SecondJudge bool
}

func NewCaptainAmericaStatus(instrument string, granularity enum.Granularity, line enum.Line, setupPrice float64, setupStatus, tradeStatus bool) *CaptainAmericaStatus {
	return &CaptainAmericaStatus{
		Instrument:  instrument,
		Granularity: granularity,
		Line:        line,
		SetupPrice:  setupPrice,
		SetupStatus: setupStatus,
		TradeStatus: tradeStatus,
		SecondJudge: false,
	}
}
