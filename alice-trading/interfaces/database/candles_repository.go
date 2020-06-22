package database

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/jinzhu/gorm"
)

// 足データのRepository
type CandlesRepository struct{}

func (rep *CandlesRepository) FindByID(db *gorm.DB, id int) (candle domain.BidAskCandles, error error) {
	candle = domain.BidAskCandles{}
	db.First(&candle, id)
	return candle, nil
}

func (rep *CandlesRepository) FindAll(db *gorm.DB) (candleList []domain.BidAskCandles) {
	candleList = []domain.BidAskCandles{}
	db.Find(&candleList)
	return candleList
}

func (rep *CandlesRepository) Create(db *gorm.DB, candle *domain.BidAskCandles) {
	db.Create(&candle)
}

func (rep *CandlesRepository) Delete(db *gorm.DB, candle *domain.BidAskCandles) {
	db.Delete(&candle)
}
