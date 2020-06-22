package database

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"github.com/jinzhu/gorm"
)

// 注文のRepository
type OrdersRepository struct{}

func (rep OrdersRepository) FindAll(db *gorm.DB) ([]domain.Orders, error) {
	var result []domain.Orders
	db.Find(&result)
	return result, nil
}

func (rep OrdersRepository) FindByID(db *gorm.DB, id int) (domain.Orders, error) {
	var result domain.Orders
	db.First(&result, id)
	return result, nil
}

func (rep OrdersRepository) FindByType(db *gorm.DB, orderType enum.Order) ([]domain.Orders, error) {
	var result []domain.Orders
	db.Where("type = ?", orderType).Find(&result)
	return result, nil
}

func (rep OrdersRepository) Create(db *gorm.DB, order *domain.Orders) {
	db.Create(&order)
}

func (rep OrdersRepository) UpdateDistance(db *gorm.DB, order *domain.Orders, param map[string]interface{}) {
	db.Model(&order).Update(param)
}
