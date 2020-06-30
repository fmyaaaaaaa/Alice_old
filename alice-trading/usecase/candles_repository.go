package usecase

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"github.com/jinzhu/gorm"
)

// 足データのRepository
type CandlesRepository interface {
	FindByID(db *gorm.DB, id int) (candle domain.BidAskCandles, error error)
	FindByInstrumentAndGranularity(db *gorm.DB, instrument string, granularity enum.Granularity) []domain.BidAskCandles
	FindAll(db *gorm.DB) (candleList []domain.BidAskCandles)
	Create(db *gorm.DB, candle *domain.BidAskCandles)
	BulkCreate(db *gorm.DB, candles *[]domain.BidAskCandles)
	Delete(db *gorm.DB, candle *domain.BidAskCandles)
}
