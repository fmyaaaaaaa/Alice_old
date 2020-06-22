package database

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/jinzhu/gorm"
)

// 紐付けのRepository
type OrderTradeBindsRepository struct{}

func (rep OrderTradeBindsRepository) FindAll(db *gorm.DB) ([]domain.OrderTradeBinds, error) {
	var result []domain.OrderTradeBinds
	db.Find(&result)
	return result, nil
}

func (rep OrderTradeBindsRepository) FindByTradeID(db *gorm.DB, tradeID int) (domain.OrderTradeBinds, error) {
	var result domain.OrderTradeBinds
	db.Where("trade_id = ? AND is_delete ", tradeID, false).Last(&result)
	return result, nil
}

func (rep OrderTradeBindsRepository) Create(db *gorm.DB, bind *domain.OrderTradeBinds) {
	db.Create(&bind)
}

func (rep OrderTradeBindsRepository) Update(db *gorm.DB, bind *domain.OrderTradeBinds, param map[string]interface{}) {
	db.Model(&bind).Updates(param)
}

func (rep OrderTradeBindsRepository) Delete(db *gorm.DB, bind *domain.OrderTradeBinds) {
	tx := db.Begin()
	var target domain.Trades
	tx.Find(&target, bind.TradeID)
	tx.Model(&target).Update("is_delete", true)
	tx.Commit()
}
