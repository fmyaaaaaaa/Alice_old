package usecase

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/jinzhu/gorm"
)

// 資金管理のRepository
type BalanceManagementsRepository interface {
	FindByInstrumentOrderByCreatedAt(db *gorm.DB, instrument string) domain.BalanceManagements
	Create(db *gorm.DB, balanceManagement *domain.BalanceManagements)
}
