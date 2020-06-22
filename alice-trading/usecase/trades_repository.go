package usecase

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/jinzhu/gorm"
)

// 取引のRepository
type TradesRepository interface {
	FindAll(db *gorm.DB) ([]domain.Trades, error)
	FindByID(db *gorm.DB, id int) (domain.Trades, error)
	Create(db *gorm.DB, trade *domain.Trades)
	Update(db *gorm.DB, trade *domain.Trades, param map[string]interface{})
}
