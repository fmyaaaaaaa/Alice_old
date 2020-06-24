package database

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"github.com/jinzhu/gorm"
)

// キャプテン・アメリカのセットアップステータスRepository
type CaptainAmericaStatusRepository struct{}

func (rep CaptainAmericaStatusRepository) FindByInstrumentAndGranularity(db *gorm.DB, instrument string, granularity enum.Granularity) domain.CaptainAmericaStatus {
	var captainAmericaStatus domain.CaptainAmericaStatus
	db.Where("instrument = ? AND granularity = ?", instrument, granularity).Find(&captainAmericaStatus)
	return captainAmericaStatus
}

func (rep CaptainAmericaStatusRepository) Create(db *gorm.DB, captainAmericaStatus *domain.CaptainAmericaStatus) {
	db.Create(&captainAmericaStatus)
}

func (rep CaptainAmericaStatusRepository) Update(db *gorm.DB, captainAmericaStatus *domain.CaptainAmericaStatus, params map[string]interface{}) {
	db.Model(&captainAmericaStatus).Updates(params)
}
