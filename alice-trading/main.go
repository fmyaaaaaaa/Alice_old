package main

import (
	"flag"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/instruments"
	"github.com/fmyaaaaaaa/Alice/alice-trading/infrastructure/cache"
	"github.com/fmyaaaaaaa/Alice/alice-trading/infrastructure/config"
	database2 "github.com/fmyaaaaaaa/Alice/alice-trading/infrastructure/database"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/api/oanda"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/database"
	"github.com/fmyaaaaaaa/Alice/alice-trading/usecase"
	"github.com/fmyaaaaaaa/Alice/alice-trading/usecase/dto"
	"sync"
	"time"
)

var cacheManager cache.AliceCacheManager
var DB *database2.DB
var instrumentsInteractor usecase.InstrumentsInteractor
var candlesInteractor usecase.CandlesInteractor

func init() {
	// configの初期化
	flag.Parse()
	if !config.InitInstance(flag.Arg(0), flag.Args()) {
		panic("failed application initialize")
	}

	// interfaceの初期化
	DB = database2.NewDB()
	instrumentsInteractor = usecase.InstrumentsInteractor{
		DB:          &database.DBRepository{DB: DB},
		Instruments: &database.InstrumentsRepository{},
	}
	candlesInteractor = usecase.CandlesInteractor{
		DB:      &database.DBRepository{DB: DB},
		Candles: &database.CandlesRepository{},
		Api:     oanda.NewCandlesApi(),
	}

	// キャッシュの構築
	cacheManager = cache.GetCacheManager()
	instrumentsInteractor.LoadInstruments()
}

func main() {
	// 銘柄ごとにgoroutineを生成し、売買を開始する。
	data := cacheManager.Get("instruments")
	instrumentList := data.([]instruments.Instruments)

	var wg sync.WaitGroup
	wg.Add(len(instrumentList))
	for _, instrument := range instrumentList {
		go func(instrumentName string) {
			defer wg.Done()
			startTrading(instrumentName)
			return
		}(instrument.Name)
	}
	wg.Wait()
}

// 指定した銘柄の取引を開始します。
func startTrading(instrumentName string) {
	// TODO:実運用で利用する時間を指定する。
	// 1分ごとに実行する
	tickPerOneMin := time.NewTicker(1 * time.Minute)
	// 1時間ごとに実行する
	tickPerOneHour := time.NewTicker(1 * time.Hour)
	// 1日ごとに実行する
	tickPerOneDay := time.NewTicker(24 * time.Hour)
	// 12時間ごとに実行する（日足の売買ルールでセットアップ発生時のみ処理を行う）
	tickPerHalfDay := time.NewTicker(12 * time.Hour)
	for {
		select {
		case <-tickPerOneMin.C:
			candlesInteractor.GetCandle(dto.NewCandlesGetDto(instrumentName, 1, enum.M1))
		case <-tickPerOneHour.C:
			candlesInteractor.GetCandle(dto.NewCandlesGetDto(instrumentName, 1, enum.H1))
		case <-tickPerOneDay.C:
			candlesInteractor.GetCandle(dto.NewCandlesGetDto(instrumentName, 1, enum.D))
		case <-tickPerHalfDay.C:
			candlesInteractor.GetCandle(dto.NewCandlesGetDto(instrumentName, 1, enum.H12))
		}
	}
}
