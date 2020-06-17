package database

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/instruments"
	"github.com/jinzhu/gorm"
)

type InstrumentsRepository struct{}

func (rep *InstrumentsRepository) FindByID(db *gorm.DB, id int) (instrument instruments.Instruments, err error) {
	instrument = instruments.Instruments{}
	db.First(&instrument, id)
	return instrument, nil
}

func (rep *InstrumentsRepository) FindAll(db *gorm.DB) (instrumentList []instruments.Instruments, err error) {
	instrumentList = []instruments.Instruments{}
	db.Find(&instrumentList)
	return instrumentList, nil
}
