package database

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"github.com/jinzhu/gorm"
)

// バックテストリザルトのRepository
type BackTestResultRepository struct{}

func (rep BackTestResultRepository) FindByInstrumentAndTradeRule(db *gorm.DB, tradeRule enum.TradeRule, instrument string) domain.BackTestResults {
	var backTestResult domain.BackTestResults
	db.Where("trade_rule = ? AND instrument = ?", tradeRule, instrument).Find(&backTestResult)
	return backTestResult
}
