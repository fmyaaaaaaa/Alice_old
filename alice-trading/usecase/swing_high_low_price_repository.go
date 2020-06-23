package usecase

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/jinzhu/gorm"
)

// スイングの高値/安値のRepository
type SwingHighLowPriceRepository interface {
	FindBySwingID(db *gorm.DB, swingID int) domain.SwingHighLowPrice
	Create(db *gorm.DB, highLowPrice *domain.SwingHighLowPrice)
	Update(db *gorm.DB, highLowPrice *domain.SwingHighLowPrice, params map[string]interface{})
}
