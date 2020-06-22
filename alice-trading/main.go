package main

import (
	"flag"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
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
var DBRepository database.DBRepository
var instrumentsInteractor usecase.InstrumentsInteractor
var candlesInteractor usecase.CandlesInteractor
var orderManager usecase.OrderManager

func init() {
	// configの初期化
	flag.Parse()
	if !config.InitInstance(flag.Arg(0), flag.Args()) {
		panic("failed application initialize")
	}

	// interfaceの初期化
	DBRepository = database.DBRepository{DB: database2.NewDB()}
	instrumentsInteractor = usecase.InstrumentsInteractor{
		DB:          &DBRepository,
		Instruments: &database.InstrumentsRepository{},
	}
	candlesInteractor = usecase.CandlesInteractor{
		DB:      &DBRepository,
		Candles: &database.CandlesRepository{},
		Api:     oanda.NewCandlesApi(),
	}
	orderManager = usecase.OrderManager{
		DB:              &DBRepository,
		Orders:          &database.OrdersRepository{},
		Trades:          &database.TradesRepository{},
		OrderTradeBinds: &database.OrderTradeBindsRepository{},
		OrdersApi:       oanda.NewOrdersApi(),
	}

	// キャッシュの構築
	cacheManager = cache.GetCacheManager()
	instrumentsInteractor.LoadInstruments()
}

func main() {
	// 銘柄ごとにgoroutineを生成し、売買を開始する。
	data := cacheManager.Get("instruments")
	instrumentList := data.([]domain.Instruments)

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
			// TODO:売買ルールと資金管理を実行する関数を実行する。
			candlesInteractor.GetCandle(dto.NewCandlesGetDto(instrumentName, 1, enum.M1))
			orderManager.DoNewMarketOrder(instrumentName, "1", "0.1")
		case <-tickPerOneHour.C:
			candlesInteractor.GetCandle(dto.NewCandlesGetDto(instrumentName, 1, enum.H1))
		case <-tickPerOneDay.C:
			candlesInteractor.GetCandle(dto.NewCandlesGetDto(instrumentName, 1, enum.D))
		case <-tickPerHalfDay.C:
			candlesInteractor.GetCandle(dto.NewCandlesGetDto(instrumentName, 1, enum.H12))
		}
	}
}
