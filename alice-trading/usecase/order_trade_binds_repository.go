package usecase

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/jinzhu/gorm"
)

// 紐付けのRepository
type OrderTradeBindsRepository interface {
	FindAll(db *gorm.DB) ([]domain.OrderTradeBinds, error)
	FindByTradeID(db *gorm.DB, tradeID int) (domain.OrderTradeBinds, error)
	Create(db *gorm.DB, trade *domain.OrderTradeBinds)
	Update(db *gorm.DB, trade *domain.OrderTradeBinds, param map[string]interface{})
	Delete(db *gorm.DB, trade *domain.OrderTradeBinds)
}
