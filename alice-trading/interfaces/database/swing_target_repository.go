package database

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"github.com/jinzhu/gorm"
)

// セットアップの検証対象となる高値/安値のRepository
type SwingTargetRepository struct{}

func (rep SwingTargetRepository) FindByID(db *gorm.DB, id int) domain.SwingTarget {
	var swingTarget domain.SwingTarget
	db.Find(&swingTarget, id)
	return swingTarget
}

func (rep SwingTargetRepository) FindByInstrumentAndGranularity(db *gorm.DB, instrument string, granularity enum.Granularity) domain.SwingTarget {
	var swingTarget domain.SwingTarget
	db.Where("instrument = ? AND granularity = ?", instrument, granularity).Find(&swingTarget)
	return swingTarget
}

func (rep SwingTargetRepository) Create(db *gorm.DB, swingTarget *domain.SwingTarget) {
	db.Create(&swingTarget)
}
