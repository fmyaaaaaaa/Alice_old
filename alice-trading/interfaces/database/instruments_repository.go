package database

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/jinzhu/gorm"
)

type InstrumentsRepository struct{}

func (rep *InstrumentsRepository) FindByID(db *gorm.DB, id int) (instrument domain.Instruments, err error) {
	instrument = domain.Instruments{}
	db.First(&instrument, id)
	return instrument, nil
}

func (rep *InstrumentsRepository) FindAll(db *gorm.DB) (instrumentList []domain.Instruments, err error) {
	instrumentList = []domain.Instruments{}
	db.Find(&instrumentList)
	return instrumentList, nil
}

func (rep *InstrumentsRepository) FindByInstrument(db *gorm.DB, instrumentName string) domain.Instruments {
	var instrument domain.Instruments
	db.Where("instrument = ?", instrumentName).Find(&instrument)
	return instrument
}
