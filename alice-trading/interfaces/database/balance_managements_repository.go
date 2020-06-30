package database

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/jinzhu/gorm"
)

// 資金管理のRepository
type BalanceManagementsRepository struct{}

func (rep BalanceManagementsRepository) FindByInstrumentOrderByCreatedAt(db *gorm.DB, instrument string) domain.BalanceManagements {
	var balanceManagement domain.BalanceManagements
	db.Where("instrument = ?", instrument).Order("created_at desc").First(&balanceManagement)
	return balanceManagement
}

func (rep BalanceManagementsRepository) Create(db *gorm.DB, balanceManagement *domain.BalanceManagements) {
	db.Create(&balanceManagement)
}
