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

func (rep CaptainAmericaStatusRepository) Reset(db *gorm.DB, instrument string, granularity enum.Granularity) {
	var target domain.CaptainAmericaStatus
	tx := db.Begin()
	tx.Where("instrument = ? AND granularity = ?", instrument, granularity).Find(&target)
	params := map[string]interface{}{
		"setup_status": false,
		"trade_status": false,
		"second_judge": false,
	}
	tx.Model(&target).Updates(params)
	tx.Commit()
}
