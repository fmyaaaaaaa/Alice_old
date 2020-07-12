package model

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"time"
)

type MidCandle struct {
	InstrumentName string
	Granularity    enum.Granularity
	Open           float64
	Close          float64
	High           float64
	Low            float64
	Line           enum.Line
	Trend          enum.Trend
	Time           time.Time
}
