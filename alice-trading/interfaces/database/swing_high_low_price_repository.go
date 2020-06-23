package database

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/jinzhu/gorm"
	"strconv"
)

// スイングの高値/安値のRepository
type SwingHighLowPriceRepository struct{}

func (rep SwingHighLowPriceRepository) FindBySwingID(db *gorm.DB, swingID int) domain.SwingHighLowPrice {
	var highLowPrice domain.SwingHighLowPrice
	db.Where("swing_id = ?", strconv.Itoa(swingID)).Find(&highLowPrice)
	return highLowPrice
}

func (rep SwingHighLowPriceRepository) Create(db *gorm.DB, highLowPrice *domain.SwingHighLowPrice) {
	db.Create(&highLowPrice)
}

func (rep SwingHighLowPriceRepository) Update(db *gorm.DB, highLowPrice *domain.SwingHighLowPrice, params map[string]interface{}) {
	db.Model(&highLowPrice).Updates(params)
}
