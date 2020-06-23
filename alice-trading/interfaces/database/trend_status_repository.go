package database

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"github.com/jinzhu/gorm"
)

// トレンドステータスのRepository
type TrendStatusRepository struct{}

// トレンドを作成します。
func (rep TrendStatusRepository) Create(db *gorm.DB, trendStatus *domain.TrendStatus) {
	db.Create(&trendStatus)
}

// 銘柄名と足種に一致するトレンドを取得します。
func (rep TrendStatusRepository) FindByInstrumentAndGranularity(db *gorm.DB, instrument string, granularity enum.Granularity) domain.TrendStatus {
	trendStatus := domain.TrendStatus{}
	db.Where("instrument = ? AND granularity = ?", instrument, granularity).Find(&trendStatus)
	return trendStatus
}

// 引数のパラメータでトレンドを更新します。
// TrendとLastSwingIDの更新を用途とします。
func (rep TrendStatusRepository) Update(db *gorm.DB, trendStatus *domain.TrendStatus, params map[string]interface{}) {
	db.Model(trendStatus).Updates(params)
}
