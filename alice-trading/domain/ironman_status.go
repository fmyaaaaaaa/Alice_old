package domain

import "github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"

// アイアンマンのセットアップステータス
type IronManStatus struct {
	ID            int
	Instrument    string
	Granularity   enum.Granularity
	SwingTargetID int
	Trend         enum.Trend
	Status        bool
}

func NewIronManStatus(instrument string, granularity enum.Granularity, swingTargetID int, trend enum.Trend) *IronManStatus {
	return &IronManStatus{
		Instrument:    instrument,
		Granularity:   granularity,
		SwingTargetID: swingTargetID,
		Trend:         trend,
		Status:        true,
	}
}
