package usecase

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"github.com/jinzhu/gorm"
)

// バックテストリザルトのRepository
type BackTestResultRepository interface {
	FindByInstrumentAndTradeRule(db *gorm.DB, tradeRule enum.TradeRule, instrument string) domain.BackTestResults
}
