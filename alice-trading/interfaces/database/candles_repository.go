package database

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/candles"
	"github.com/jinzhu/gorm"
)

// 足データのRepository
type CandlesRepository struct{}

func (rep *CandlesRepository) FindByID(db *gorm.DB, id int) (candle candles.BidAskCandles, error error) {
	candle = candles.BidAskCandles{}
	db.First(&candle, id)
	return candle, nil
}

func (rep *CandlesRepository) FindAll(db *gorm.DB) (candleList []candles.BidAskCandles) {
	candleList = []candles.BidAskCandles{}
	db.Find(&candleList)
	return candleList
}

func (rep *CandlesRepository) Create(db *gorm.DB, candle *candles.BidAskCandles) {
	db.Create(&candle)
}

func (rep *CandlesRepository) Delete(db *gorm.DB, candle *candles.BidAskCandles) {
	db.Delete(&candle)
}
