package domain

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"time"
)

// 売買ルールのセットアップステータス
type TradeRuleStatus struct {
	ID          int
	TradeRule   enum.TradeRule
	Instrument  string
	Granularity enum.Granularity
	CandleTime  time.Time
	Status      bool
}

func NewTradeRuleStatus(tradeRule enum.TradeRule, instrument string, granularity enum.Granularity, candleTime time.Time) *TradeRuleStatus {
	return &TradeRuleStatus{
		TradeRule:   tradeRule,
		Instrument:  instrument,
		Granularity: granularity,
		CandleTime:  candleTime,
		Status:      true,
	}
}
