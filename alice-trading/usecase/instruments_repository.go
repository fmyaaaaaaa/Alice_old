package usecase

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/jinzhu/gorm"
)

// 銘柄関連のRepository
type InstrumentsRepository interface {
	FindByID(db *gorm.DB, id int) (instrument domain.Instruments, err error)
	FindAll(db *gorm.DB) (instrumentList []domain.Instruments, err error)
}
