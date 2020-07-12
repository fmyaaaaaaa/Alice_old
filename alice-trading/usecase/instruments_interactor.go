package usecase

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"github.com/fmyaaaaaaa/Alice/alice-trading/infrastructure/cache"
	"log"
)

// 銘柄関連のユースケース
type InstrumentsInteractor struct {
	DB          DBRepository
	Instruments InstrumentsRepository
}

// 銘柄をDBから取得し、キャッシュに保存します。
// アプリ起動時のみ利用します。
func (i *InstrumentsInteractor) LoadInstruments() {
	cacheManager := cache.GetCacheManager()
	instrumentList, err := i.GetAll()
	if err != nil {
		log.Fatal(err)
	}
	cacheManager.Set("instruments", instrumentList, enum.NoExpiration)
}

// 主キーをもとに銘柄を取得します。
func (i *InstrumentsInteractor) Get(id int) (instrument domain.Instruments, err error) {
	db := i.DB.Connect()
	result, err := i.Instruments.FindByID(db, id)
	if err != nil {
		return domain.Instruments{}, err
	}
	return result, nil
}

// 銘柄を全件取得します。
func (i *InstrumentsInteractor) GetAll() (instrumentList []domain.Instruments, err error) {
	db := i.DB.Connect()
	result, err := i.Instruments.FindAll(db)
	if err != nil {
		panic(err)
	}
	return result, nil
}

// 銘柄名を指定して銘柄情報を取得します。
func (i *InstrumentsInteractor) GetInstrument(instrument string) domain.Instruments {
	db := i.DB.Connect()
	return i.Instruments.FindByInstrument(db, instrument)
}
