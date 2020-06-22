package usecase

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/jinzhu/gorm"
)

// 足データのRepository
type CandlesRepository interface {
	FindByID(db *gorm.DB, id int) (candle domain.BidAskCandles, error error)
	FindAll(db *gorm.DB) (candleList []domain.BidAskCandles)
	Create(db *gorm.DB, candle *domain.BidAskCandles)
	Delete(db *gorm.DB, candle *domain.BidAskCandles)
}
