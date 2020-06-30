package database

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"github.com/jinzhu/gorm"
)

// 足データのRepository
type CandlesRepository struct{}

func (rep CandlesRepository) FindByID(db *gorm.DB, id int) (candle domain.BidAskCandles, error error) {
	candle = domain.BidAskCandles{}
	db.First(&candle, id)
	return candle, nil
}

func (rep CandlesRepository) FindByInstrumentAndGranularity(db *gorm.DB, instrument string, granularity enum.Granularity) []domain.BidAskCandles {
	var candles []domain.BidAskCandles
	db.Where("instrument_name = ? AND granularity = ?", instrument, granularity).Find(&candles)
	return candles
}

func (rep CandlesRepository) FindAll(db *gorm.DB) (candleList []domain.BidAskCandles) {
	candleList = []domain.BidAskCandles{}
	db.Find(&candleList)
	return candleList
}

func (rep CandlesRepository) Create(db *gorm.DB, candle *domain.BidAskCandles) {
	db.Create(&candle)
}

func (rep CandlesRepository) BulkCreate(db *gorm.DB, candles *[]domain.BidAskCandles) {
	tx := db.Begin()
	for _, candle := range *candles {
		db.Create(&candle)
	}
	tx.Commit()
}

func (rep CandlesRepository) Delete(db *gorm.DB, candle *domain.BidAskCandles) {
	db.Delete(&candle)
}
