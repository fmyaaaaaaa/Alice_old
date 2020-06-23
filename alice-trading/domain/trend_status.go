package domain

import "github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"

// トレンドのステータス
type TrendStatus struct {
	ID          int
	Instrument  string
	Granularity enum.Granularity
	Trend       enum.Trend
	LastSwingID int
}

func NewTrendStatus(instrument string, granularity enum.Granularity, trend enum.Trend, swingID int) *TrendStatus {
	return &TrendStatus{
		Instrument:  instrument,
		Granularity: granularity,
		Trend:       trend,
		LastSwingID: swingID,
	}
}
