package domain

import "github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"

// セットアップの検証対象となる高値/安値
type SwingTarget struct {
	ID          int
	Instrument  string
	Granularity enum.Granularity
	SwingID     int
}

func NewSwingTarget(instrument string, granularity enum.Granularity, swingID int) *SwingTarget {
	return &SwingTarget{
		Instrument:  instrument,
		Granularity: granularity,
		SwingID:     swingID,
	}
}
