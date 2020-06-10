package usecase

import "github.com/fmyaaaaaaa/Alice/alice-trading/domain/instruments"

type InstrumentsInteractor struct {
	DB          DBRepository
	Instruments InstrumentsRepository
}

func (i *InstrumentsInteractor) Get(id int) (instrument instruments.Instruments, err error) {
	db := i.DB.Connect()
	result, err := i.Instruments.FindByID(db, id)
	if err != nil {
		return instruments.Instruments{}, err
	}
	return result, nil

}
