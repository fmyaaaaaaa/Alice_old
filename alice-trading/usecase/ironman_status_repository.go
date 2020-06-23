package usecase

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"github.com/jinzhu/gorm"
)

// アイアンマンのセットアップステータスのRepository
type IronManStatusRepository interface {
	FindByInstrumentAndGranularity(db *gorm.DB, instrument string, granularity enum.Granularity) domain.IronManStatus
	Create(db *gorm.DB, ironManStatus *domain.IronManStatus)
	Update(db *gorm.DB, ironManStatus *domain.IronManStatus, params map[string]interface{})
}
