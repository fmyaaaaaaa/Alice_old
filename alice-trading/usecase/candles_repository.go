package usecase

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/candles"
	"github.com/jinzhu/gorm"
)

// 足データのRepository
type CandlesRepository interface {
	FindByID(db *gorm.DB, id int) (candle candles.BidAskCandles, error error)
	FindAll(db *gorm.DB) (candleList []candles.BidAskCandles)
	Create(db *gorm.DB, candle *candles.BidAskCandles)
	Delete(db *gorm.DB, candle *candles.BidAskCandles)
}
