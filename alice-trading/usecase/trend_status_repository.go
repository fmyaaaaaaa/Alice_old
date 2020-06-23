package usecase

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"github.com/jinzhu/gorm"
)

// トレンドステータスのRepository
type TrendStatusRepository interface {
	Create(db *gorm.DB, trendStatus *domain.TrendStatus)
	FindByInstrumentAndGranularity(db *gorm.DB, instrument string, duration enum.Granularity) domain.TrendStatus
	Update(db *gorm.DB, trendStatus *domain.TrendStatus, params map[string]interface{})
}
