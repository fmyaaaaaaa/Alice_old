package usecase

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/instruments"
	"github.com/jinzhu/gorm"
)

type InstrumentsRepository interface {
	FindByID(db *gorm.DB, id int) (instrument instruments.Instruments, err error)
}
