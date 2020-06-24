package usecase

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"github.com/jinzhu/gorm"
)

// 売買ルールのセットアップステータスのRepository
type TradeRuleStatusRepository interface {
	FindTargetByTradeRuleAndInstrumentAndGranularity(db *gorm.DB, tradeRule enum.TradeRule, instrument string, granularity enum.Granularity) domain.TradeRuleStatus
	FindByTradeRuleAndInstrumentAndGranularity(db *gorm.DB, tradeRule enum.TradeRule, instrument string, granularity enum.Granularity) domain.TradeRuleStatus
	Create(db *gorm.DB, tradeRuleStatus *domain.TradeRuleStatus)
	Update(db *gorm.DB, tradeRuleStatus *domain.TradeRuleStatus, params map[string]interface{})
}
