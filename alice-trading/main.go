package main

import (
	"flag"
	"fmt"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"github.com/fmyaaaaaaa/Alice/alice-trading/infrastructure/cache"
	"github.com/fmyaaaaaaa/Alice/alice-trading/infrastructure/config"
	database2 "github.com/fmyaaaaaaa/Alice/alice-trading/infrastructure/database"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/api/oanda"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/database"
	"github.com/fmyaaaaaaa/Alice/alice-trading/usecase"
	"github.com/fmyaaaaaaa/Alice/alice-trading/usecase/dto"
	"github.com/fmyaaaaaaa/Alice/alice-trading/usecase/money"
	"github.com/fmyaaaaaaa/Alice/alice-trading/usecase/rule"
	"log"
	"strconv"
	"sync"
	"time"
)

var cacheManager cache.AliceCacheManager
var DBRepository database.DBRepository
var instrumentsInteractor usecase.InstrumentsInteractor
var candlesInteractor usecase.CandlesInteractor
var orderManager usecase.OrderManager
var avengers rule.Avengers
var ironMan rule.IronMan
var captainAmerica rule.CaptainAmerica
var balanceManager money.BalanceManager
var accountManager money.AccountManager

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
	avengers = rule.Avengers{
		DB:                &DBRepository,
		TrendStatus:       &database.TrendStatusRepository{},
		Sequence:          &database.SequenceRepository{},
		TradeRuleStatus:   &database.TradeRuleStatusRepository{},
		SwingHighLowPrice: &database.SwingHighLowPriceRepository{},
		SwingTarget:       &database.SwingTargetRepository{},
	}
	ironMan = rule.IronMan{
		DB:                &DBRepository,
		SwingHighLowPrice: &database.SwingHighLowPriceRepository{},
		SwingTarget:       &database.SwingTargetRepository{},
		IronManStatus:     &database.IronManStatusRepository{},
		TradeRuleStatus:   &database.TradeRuleStatusRepository{},
	}
	captainAmerica = rule.CaptainAmerica{
		DB:                   &DBRepository,
		TrendStatus:          &database.TrendStatusRepository{},
		CaptainAmericaStatus: &database.CaptainAmericaStatusRepository{},
		TradeRuleStatus:      &database.TradeRuleStatusRepository{},
	}
	balanceManager = money.BalanceManager{
		DB:                 &DBRepository,
		PricesApi:          oanda.NewPricesApi(),
		BackTestResult:     &database.BackTestResultRepository{},
		BalanceManagements: &database.BalanceManagementsRepository{},
	}
	accountManager = money.AccountManager{
		DB:           &DBRepository,
		Accounts:     &database.AccountsRepository{},
		Trades:       &database.TradesRepository{},
		Positions:    &database.PositionsRepository{},
		AccountsApi:  oanda.NewAccountApi(),
		TradesApi:    oanda.NewTradeApi(),
		PositionsApi: oanda.NewPositionsApi(),
	}

	// キャッシュの構築
	cacheManager = cache.GetCacheManager()
	instrumentsInteractor.LoadInstruments()
}

func main() {
	data := cacheManager.Get("instruments")
	instruments := data.([]domain.Instruments)

	log.Println("Starting create candles cache")
	// 足データをキャッシュに構築する。
	for _, instrument := range instruments {
		if ok := candlesInteractor.InitializeCandle(instrument, enum.H1); ok {
			judgementInitialCandles(instrument, enum.H1)
		}
		if ok := candlesInteractor.InitializeCandle(instrument, enum.D); ok {
			judgementInitialCandles(instrument, enum.D)
		}
	}

	log.Println("Starting Application")
	// 銘柄ごとにgoroutineを生成し、売買を開始する。
	var wg sync.WaitGroup
	wg.Add(len(instruments) + 1)
	for _, instrument := range instruments {
		go func(instrument domain.Instruments) {
			defer wg.Done()
			startTrading(instrument)
			return
		}(instrument)
	}
	go startWatchAccountInformation()
	wg.Wait()
}

// 指定した銘柄の取引を開始します。
func startTrading(instrument domain.Instruments) {
	// 1分ごとに実行する
	tickPerOneMin := time.NewTicker(1 * time.Minute)
	// 1時間ごとに実行する
	tickPerOneHour := time.NewTicker(1 * time.Hour)
	// 1日ごとに実行する
	tickPerOneDay := time.NewTicker(24 * time.Hour)
	// 12時間ごとに実行する（日足の売買ルールでセットアップ発生時のみ処理を行う）
	tickPerHalfDay := time.NewTicker(12 * time.Minute)

	for {
		select {
		case <-tickPerOneMin.C:
			// ポジションを保持している場合、ポジション情報を更新する。
			if accountManager.HasPosition(instrument.Instrument) {
				accountManager.UpdatePositionInformation(instrument.Instrument)
				ok, currentAccountLevel, tradeRule := balanceManager.JudgementProfit(accountManager.GetPosition(instrument.Instrument))
				if ok {
					distance := balanceManager.GetTradingDistance(instrument)
					units := strconv.FormatFloat(accountManager.GetPosition(instrument.Instrument).Units, 'f', 0, 64)
					trade := orderManager.DoNewMarketOrder(instrument.Instrument, units, strconv.FormatFloat(distance, 'f', 5, 64))
					position := accountManager.CreateOrUpdatePosition(instrument.Instrument)
					balanceManager.RegisterNextBalanceManagements(instrument, trade, position, tradeRule, currentAccountLevel, distance)
					log.Print("Complete New Order :", instrument.Instrument)
				}
			}
		case <-tickPerOneHour.C:
			//candle := candlesInteractor.GetCandle(dto.NewCandlesGetDto(instrument.Instrument, 1, enum.H1))[0]
			//judgementCurrentCandle(&candle, &instrument, enum.H1)
			//// セットアップの検証
			//ironMan.JudgementSetup(&candle, instrument.Instrument, enum.H1)
			//
			//// トレード計画の検証
			//tradeRuleStatus, ok := avengers.IsExistSetUpTradeRule(enum.IronMan, instrument.Instrument, enum.H1)
			//if ok {
			//	// トレード計画を検証する。
			//	isOrder, units := ironMan.JudgementTradePlan(tradeRuleStatus, &candle, instrument.Instrument, enum.H1)
			//	// ポジション状況を取得し、注文可能であれば注文を実行する。
			//	if isOrder {
			//		if hasPosition := accountManager.HasPosition(instrument.Instrument); !hasPosition {
			//			distance := balanceManager.GetTradingDistance(instrument)
			//			trade := orderManager.DoNewMarketOrder(instrument.Instrument, units, strconv.FormatFloat(distance, 'f', 5, 64))
			//			accountManager.CreateOrUpdatePosition(instrument.Instrument)
			//			balanceManager.RegisterFirstBalanceManagements(instrument, trade, enum.IronMan, distance)
			//			log.Print("Complete New Order :", instrument.Instrument)
			//		}
			//	}
			//}
		case <-tickPerOneDay.C:
			log.Println("Start Judgement Setup CaptainAmerica 24 hours :", instrument.Instrument)
			candle := candlesInteractor.GetCandle(dto.NewCandlesGetDto(instrument.Instrument, 1, enum.H1))[0]
			candles := judgementCurrentCandle(&candle, &instrument, enum.H1)
			// セットアップの検証
			captainAmerica.JudgementSetup(&candles[len(candles)-2], &candles[len(candles)-1], instrument.Instrument, enum.H1)
			tradeStatus, ok := avengers.IsExistSetUpTradeRule(enum.CaptainAmerica, instrument.Instrument, enum.H1)
			if ok && captainAmerica.IsExistSecondJudgementTradePlan(instrument.Instrument, enum.H1) {
				isOrder, units := captainAmerica.JudgementTradePlan(tradeStatus, &candles[len(candles)-1], instrument.Instrument, enum.H1)
				if isOrder && accountManager.HasPosition(instrument.Instrument) && canCreateNewOrder() {
					// 以下の条件時に注文可能とする。
					// ①トレード計画検証結果がOK　②同じ銘柄でポジション未保持　③オープンポジションが６つ以下
					distance := balanceManager.GetTradingDistance(instrument)
					trade := orderManager.DoNewMarketOrder(instrument.Instrument, units, strconv.FormatFloat(distance, 'f', 5, 64))
					accountManager.CreateOrUpdatePosition(instrument.Instrument)
					balanceManager.RegisterFirstBalanceManagements(instrument, trade, enum.CaptainAmerica, distance)
					log.Print("Complete New Order :", instrument.Instrument)
				}
			}
		case <-tickPerHalfDay.C:
			candle := candlesInteractor.GetCandle(dto.NewCandlesGetDto(instrument.Instrument, 1, enum.M30))[0]
			const TimeFormat = "15:04:05"
			targetTime := candle.Candles.Time.Format(TimeFormat)
			if targetTime == "12:00:00" {
				log.Println("Start Judgement TradePlan CaptainAmerica 12 hours :", instrument.Instrument)
				tradeStatus, ok := avengers.IsExistSetUpTradeRule(enum.CaptainAmerica, instrument.Instrument, enum.H1)
				if ok {
					isOrder, units := captainAmerica.JudgementTradePlan(tradeStatus, &candle, instrument.Instrument, enum.H1)
					// 以下の条件時に注文可能とする。
					// ①トレード計画検証結果がOK　②同じ銘柄でポジション未保持　③オープンポジションが６つ以下
					if isOrder && !accountManager.HasPosition(instrument.Instrument) && canCreateNewOrder() {
						distance := balanceManager.GetTradingDistance(instrument)
						trade := orderManager.DoNewMarketOrder(instrument.Instrument, units, strconv.FormatFloat(distance, 'f', 5, 64))
						accountManager.CreateOrUpdatePosition(instrument.Instrument)
						balanceManager.RegisterFirstBalanceManagements(instrument, trade, enum.CaptainAmerica, distance)
						log.Print("Complete New Order :", instrument.Instrument)
					}
				}
			}
		}
	}
}

// アカウント情報を1分間隔で更新します。
func startWatchAccountInformation() {
	tickPerOneMin := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-tickPerOneMin.C:
			accountManager.UpdateAccountInformation()
		}
	}
}

// 現在の足データに関する線種、トレンドを判定します。
func judgementCurrentCandle(currentCandle *domain.BidAskCandles, instrument *domain.Instruments, granularity enum.Granularity) []domain.BidAskCandles {
	candles := getCandles(instrument.Instrument, granularity)
	candles = append(candles, *currentCandle)

	avengers.JudgementLine(&candles[len(candles)-1])
	avengers.JudgementTrend(&candles[len(candles)-4], &candles[len(candles)-3], &candles[len(candles)-2], &candles[len(candles)-1])
	avengers.JudgementSwingAndAllTrend(&candles[len(candles)-1], instrument.Instrument, granularity)

	setCandles(candles, instrument.Instrument, granularity)
	candlesInteractor.Create(&candles[len(candles)-1])
	return candles
}

// APIで取得した足データに対して、線種、トレンドを判定し付与します。
func judgementInitialCandles(instrument domain.Instruments, granularity enum.Granularity) {
	candles := getCandles(instrument.Instrument, granularity)
	for i := range candles {
		avengers.JudgementLine(&candles[i])
		candles[i].Trend = enum.Range
	}
	avengers.JudgementTrend(&candles[0], &candles[1], &candles[2], &candles[3])
	avengers.CreateTrendStatus(domain.NewTrendStatus(instrument.Instrument, granularity, candles[3].Trend, avengers.Increment(enum.Swing)))
	candlesInteractor.CreateBulkCandles(candles)
	setCandles(candles, instrument.Instrument, granularity)
}

// オープンポジションが6つ以下の場合のみ新規注文を受け付けます。
func canCreateNewOrder() bool {
	account := accountManager.GetAccount()
	if account.OpenPositionCount <= 6 {
		return true
	}
	return false
}

// キャッシュから足データを取得します。
// 足データのKeyは candles-USD_JPY-H1
func getCandles(instrument string, granularity enum.Granularity) []domain.BidAskCandles {
	cacheName := fmt.Sprintf("candles-%s-%s", instrument, granularity)
	data := cacheManager.Get(cacheName)
	return data.([]domain.BidAskCandles)
}

// キャッシュに足データを保存します。
// 足データのKeyは candles-USD_JPY-H1
func setCandles(candles []domain.BidAskCandles, instrument string, granularity enum.Granularity) {
	cacheName := fmt.Sprintf("candles-%s-%s", instrument, granularity)
	cacheManager.Set(cacheName, candles, enum.NoExpiration)
}
