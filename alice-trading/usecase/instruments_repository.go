package usecase

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/instruments"
	"github.com/jinzhu/gorm"
)

// 銘柄関連のRepository
type InstrumentsRepository interface {
	FindByID(db *gorm.DB, id int) (instrument instruments.Instruments, err error)
	FindAll(db *gorm.DB) (instrumentList []instruments.Instruments, err error)
}
