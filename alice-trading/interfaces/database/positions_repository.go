package database

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/jinzhu/gorm"
)

// ポジションのRepository
type PositionsRepository struct{}

func (rep PositionsRepository) FindByInstrument(db *gorm.DB, instrument string) domain.Positions {
	var position domain.Positions
	db.Where("instrument = ?", instrument).Find(&position)
	return position
}

func (rep PositionsRepository) FindAll(db *gorm.DB) []domain.Positions {
	var positions []domain.Positions
	db.Find(&positions)
	return positions
}

func (rep PositionsRepository) CreateOrUpdate(db *gorm.DB, position *domain.Positions) {
	tx := db.Begin()
	var target domain.Positions
	tx.Where("instrument = ?", position.Instrument).Find(&target)
	if target.ID != 0 {
		position.ID = target.ID
		tx.Save(position)
	} else {
		tx.Create(position)
	}
	tx.Commit()
}

func (rep PositionsRepository) Update(db *gorm.DB, position *domain.Positions) {
	tx := db.Begin()
	var target domain.Positions
	tx.Where("instrument = ?", position.Instrument).Find(&target)
	position.ID = target.ID
	tx.Save(&position)
	tx.Commit()
}

func (rep PositionsRepository) BulkUpdate(db *gorm.DB, positions *[]domain.Positions) {
	tx := db.Begin()
	for _, position := range *positions {
		var target domain.Positions
		tx.Where("instrument = ?", position.Instrument).Find(&target)
		position.ID = target.ID
		tx.Save(&position)
	}
	tx.Commit()
}
