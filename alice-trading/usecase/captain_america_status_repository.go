package usecase

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"github.com/jinzhu/gorm"
)

// キャプテン・アメリカのセットアップステータスRepository
type CaptainAmericaStatusRepository interface {
	FindByInstrumentAndGranularity(db *gorm.DB, instrument string, granularity enum.Granularity) domain.CaptainAmericaStatus
	Create(db *gorm.DB, captainAmericaStatus *domain.CaptainAmericaStatus)
	Update(db *gorm.DB, captainAmericaStatus *domain.CaptainAmericaStatus, params map[string]interface{})
}
