package usecase

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"github.com/jinzhu/gorm"
)

// 注文のRepository
type OrdersRepository interface {
	FindAll(db *gorm.DB) ([]domain.Orders, error)
	FindByID(db *gorm.DB, id int) (domain.Orders, error)
	FindLastByInstrumentAndOrder(db *gorm.DB, instrument string, order enum.Order) (domain.Orders, error)
	FindByType(db *gorm.DB, orderType enum.Order) ([]domain.Orders, error)
	Create(db *gorm.DB, order *domain.Orders)
	UpdateDistance(db *gorm.DB, order *domain.Orders, param map[string]interface{})
}
