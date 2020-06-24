package database

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"github.com/jinzhu/gorm"
)

// 売買ルールのセットアップステータスのRepository
type TradeRuleStatusRepository struct{}

func (rep TradeRuleStatusRepository) FindTargetByTradeRuleAndInstrumentAndGranularity(db *gorm.DB, tradeRule enum.TradeRule, instrument string, granularity enum.Granularity) domain.TradeRuleStatus {
	var tradeRuleStatus domain.TradeRuleStatus
	db.Where("trade_rule = ? AND instrument = ? AND granularity = ? AND status = ?", tradeRule, instrument, granularity, true).Find(&tradeRuleStatus)
	return tradeRuleStatus
}

func (rep TradeRuleStatusRepository) FindByTradeRuleAndInstrumentAndGranularity(db *gorm.DB, tradeRule enum.TradeRule, instrument string, granularity enum.Granularity) domain.TradeRuleStatus {
	var tradeRuleStatus domain.TradeRuleStatus
	db.Where("trade_rule = ? AND instrument = ? AND granularity = ?", tradeRule, instrument, granularity).Find(&tradeRuleStatus)
	return tradeRuleStatus
}

func (rep TradeRuleStatusRepository) Create(db *gorm.DB, tradeRuleStatus *domain.TradeRuleStatus) {
	db.Create(&tradeRuleStatus)
}

func (rep TradeRuleStatusRepository) Update(db *gorm.DB, tradeRuleStatus *domain.TradeRuleStatus, params map[string]interface{}) {
	db.Model(&tradeRuleStatus).Updates(params)
}
