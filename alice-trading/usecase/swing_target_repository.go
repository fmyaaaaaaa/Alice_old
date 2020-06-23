package usecase

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"github.com/jinzhu/gorm"
)

// セットアップの検証対象となる高値/安値のRepository
type SwingTargetRepository interface {
	FindByID(db *gorm.DB, id int) domain.SwingTarget
	FindByInstrumentAndGranularity(db *gorm.DB, instrument string, granularity enum.Granularity) domain.SwingTarget
	Create(db *gorm.DB, swingTarget *domain.SwingTarget)
}
