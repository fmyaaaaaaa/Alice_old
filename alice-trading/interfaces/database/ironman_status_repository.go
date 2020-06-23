package database

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"github.com/jinzhu/gorm"
)

// アイアンマンのセットアップステータスのRepository
type IronManStatusRepository struct{}

func (rep IronManStatusRepository) FindByInstrumentAndGranularity(db *gorm.DB, instrument string, granularity enum.Granularity) domain.IronManStatus {
	var ironManStatus domain.IronManStatus
	db.Where("instrument = ? AND granularity = ?", instrument, granularity).Last(&ironManStatus)
	return ironManStatus
}

func (rep IronManStatusRepository) Create(db *gorm.DB, ironManStatus *domain.IronManStatus) {
	db.Create(&ironManStatus)
}

func (rep IronManStatusRepository) Update(db *gorm.DB, ironManStatus *domain.IronManStatus, params map[string]interface{}) {
	db.Model(&ironManStatus).Updates(params)
}
