package domain

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"time"
)

// TODO:AccountLevelの名称を検討する。
// 資金管理
type BalanceManagements struct {
	ID                  int
	TradeID             int
	Instrument          string
	TradeRule           enum.TradeRule
	CurrentAccountLevel float64
	NextAccountLevel    float64
	PositionSize        float64
	ExecPrice           float64
	Distance            float64
	Delta               float64
	CreatedAt           time.Time
}

func NewBalanceManagements(tradeID int, instrument string, tradeRule enum.TradeRule, currentAccountLevel, nextAccountLevel, positionSize, execPrice, distance, delta float64) *BalanceManagements {
	return &BalanceManagements{
		TradeID:             tradeID,
		Instrument:          instrument,
		TradeRule:           tradeRule,
		CurrentAccountLevel: currentAccountLevel,
		NextAccountLevel:    nextAccountLevel,
		PositionSize:        positionSize,
		ExecPrice:           execPrice,
		Distance:            distance,
		Delta:               delta,
		CreatedAt:           time.Now(),
	}
}
