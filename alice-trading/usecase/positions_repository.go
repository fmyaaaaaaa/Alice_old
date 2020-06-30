package usecase

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/jinzhu/gorm"
)

// ポジションのRepository
type PositionsRepository interface {
	FindByInstrument(db *gorm.DB, instrument string) domain.Positions
	FindAll(db *gorm.DB) []domain.Positions
	CreateOrUpdate(db *gorm.DB, position *domain.Positions)
	Update(db *gorm.DB, position *domain.Positions)
	BulkUpdate(db *gorm.DB, positions *[]domain.Positions)
}
