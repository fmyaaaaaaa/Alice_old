package database

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/jinzhu/gorm"
)

// 取引のRepository
type TradesRepository struct{}

func (rep TradesRepository) FindAll(db *gorm.DB) ([]domain.Trades, error) {
	var result []domain.Trades
	db.Find(&result)
	return result, nil
}

func (rep TradesRepository) FindByID(db *gorm.DB, id int) (domain.Trades, error) {
	var result domain.Trades
	db.First(&result, id)
	return result, nil
}

func (rep TradesRepository) Create(db *gorm.DB, trade *domain.Trades) {
	db.Create(&trade)
}

func (rep TradesRepository) Update(db *gorm.DB, trade *domain.Trades, param map[string]interface{}) {
	db.Model(&trade).Updates(param)
}
