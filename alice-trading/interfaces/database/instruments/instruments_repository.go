package instruments

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/instruments"
	"github.com/jinzhu/gorm"
)

type Repository struct{}

func (rep *Repository) FindByID(db *gorm.DB, id int) (instrument instruments.Instruments, err error) {
	instrument = instruments.Instruments{}
	db.First(&instrument, id)
	return instrument, nil
}
