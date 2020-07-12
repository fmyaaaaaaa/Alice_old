package main

import (
	"flag"
	"fmt"
	"github.com/fmyaaaaaaa/Alice/alice-trading/backtest/controller"
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
	// backTestモードの場合
	if flag.Arg(0) == "backTest" {
		controller.StartServer()
		return
	}

	data := cacheManager.Get("instruments")
	instruments := data.([]domain.Instruments)

	// 足データをキャッシュに構築する。
	log.Println("Starting create candles cache")
	createCandleData(instruments, enum.TradeType(config.GetInstance().Rule.TradeType))

	// 銘柄ごとにgoroutineを生成し、売買を開始する。
	log.Println("Starting Application")
	var wg sync.WaitGroup
	wg.Add(len(instruments) + 1)
	for _, instrument := range instruments {
		go func(instrument domain.Instruments) {
			defer wg.Done()
			handleTradeRule(instrument)
			return
		}(instrument)
	}
	go startWatchAccountInformation()
	wg.Wait()
}

// 引数の足種で足データを取得またはロードします。
// 必要に応じて足データの線種、トレンドを判定します。
func createCandleData(instruments []domain.Instruments, tradeType enum.TradeType) {
	// 足データをキャッシュに構築する。
	var granularity enum.Granularity
	switch tradeType {
	case enum.DayTrade:
		granularity = enum.H1
	case enum.SwingTrade:
		granularity = enum.D
	}

	for _, instrument := range instruments {
		if ok := candlesInteractor.InitializeCandle(instrument, granularity); ok {
			candles := judgementInitialCandles(instrument, granularity)
			captainAmerica.JudgementSetup(&candles[2], &candles[3], instrument.Instrument, granularity)
		}
	}
}

// 売買ルールによるハンドリングを行います。
func handleTradeRule(instrument domain.Instruments) {
	switch enum.TradeRule(config.GetInstance().Rule.TradeRule) {
	case enum.CaptainAmerica:
		switch enum.TradeType(config.GetInstance().Rule.TradeType) {
		case enum.DayTrade:
			startDayTradingCaptainAmerica(instrument)
		case enum.SwingTrade:
			startSwingTradingCaptainAmerica(instrument)
		}
	case enum.IronMan:
		switch enum.TradeType(config.GetInstance().Rule.TradeType) {
		case enum.DayTrade:
			// TODO 正式実装する
			startDayTradingIronMan()
		case enum.SwingTrade:
			// TODO 正式実装する
			startSwingTradingIronMan()
		}
	}
}

// デイトレードのキャプテンアメリカを開始します。
func startDayTradingCaptainAmerica(instrument domain.Instruments) {
	// 1分ごとに実行する
	tickPerOneMin := time.NewTicker(1 * time.Minute)
	// 1時間ごとに実行する
	tickPerOneHour := time.NewTicker(1 * time.Hour)
	// 30分ごとに実行する（売買ルールでセットアップ発生時のみ処理を行う）
	tickPerHalfHour := time.NewTicker(30 * time.Minute)

	for {
		select {
		case <-tickPerOneMin.C:
			if ok, _ := accountManager.HasPosition(instrument.Instrument); ok {
				if ok, position := accountManager.UpdatePositionInformation(instrument.Instrument); ok {
					handlePosition(position, instrument.Instrument, enum.H1)
					if ok, currentAccountLevel, tradeRule := balanceManager.JudgementProfit(accountManager.GetPosition(instrument.Instrument)); ok {
						distance, trade, position := doCreateNewOrder(instrument, strconv.FormatFloat(accountManager.GetPosition(instrument.Instrument).Units, 'f', 0, 64))
						balanceManager.RegisterNextBalanceManagements(instrument, trade, position, tradeRule, currentAccountLevel, distance)
					}
				}
			}
		case <-tickPerOneHour.C:
			log.Println("Start Judgement Setup CaptainAmerica DayTrade 1 hours :", instrument.Instrument)
			candle := candlesInteractor.GetCandle(dto.NewCandlesGetDto(instrument.Instrument, 2, enum.H1))[0]
			candles := doSetupCaptainAmerica(candle, instrument, enum.H1)
			tradeStatus, ok := avengers.IsExistSetUpTradeRule(enum.CaptainAmerica, instrument.Instrument, enum.H1)
			if ok && captainAmerica.IsExistSecondJudgementTradePlan(instrument.Instrument, enum.H1) {
				log.Println("Start Judgement TradePlan CaptainAmerica DayTrade Second Judge :", instrument.Instrument)
				isOrder, units, captainAmericaStatus := captainAmerica.JudgementTradePlan(tradeStatus, &candles[len(candles)-1], instrument.Instrument, enum.H1)
				if ok, _ := accountManager.HasPosition(instrument.Instrument); !ok && isOrder && canCreateNewOrder(instrument.Instrument) {
					distance, trade, _ := doCreateNewOrder(instrument, units)
					captainAmerica.CompleteTradeStatus(&captainAmericaStatus)
					balanceManager.RegisterFirstBalanceManagements(instrument, trade, enum.CaptainAmerica, distance)
				}
			}
		case <-tickPerHalfHour.C:
			candle := candlesInteractor.GetCandle(dto.NewCandlesGetDto(instrument.Instrument, 2, enum.M30))[0]
			const TimeFormat = "04:05"
			if candle.Candles.Time.Format(TimeFormat) == "00:00" {
				log.Println("Start Judgement TradePlan CaptainAmerica DayTrade First Judge :", instrument.Instrument)
				// セットアップの検証時には足データの足種を判定しているため
				tradeStatus, ok := avengers.IsExistSetUpTradeRule(enum.CaptainAmerica, instrument.Instrument, enum.H1)
				if ok {
					avengers.JudgementLine(&candle)
					isOrder, units, captainAmericaStatus := captainAmerica.JudgementTradePlan(tradeStatus, &candle, instrument.Instrument, enum.H1)
					if ok, _ := accountManager.HasPosition(instrument.Instrument); !ok && isOrder && canCreateNewOrder(instrument.Instrument) {
						distance, trade, _ := doCreateNewOrder(instrument, units)
						captainAmerica.CompleteTradeStatus(&captainAmericaStatus)
						balanceManager.RegisterFirstBalanceManagements(instrument, trade, enum.CaptainAmerica, distance)
					}
				}
			}
		}
	}
}

// スイングトレードのキャプテンアメリカを開始します。
func startSwingTradingCaptainAmerica(instrument domain.Instruments) {
	// 1分ごとに実行する
	tickPerOneMin := time.NewTicker(1 * time.Minute)
	// 24時間ごとに実行する
	tickPerDay := time.NewTicker(24 * time.Hour)
	// 12時間ごとに実行する（売買ルールでセットアップ発生時のみ処理を行う）
	tickPerHalfDay := time.NewTicker(12 * time.Hour)

	for {
		select {
		case <-tickPerOneMin.C:
			if ok, _ := accountManager.HasPosition(instrument.Instrument); ok {
				if ok, position := accountManager.UpdatePositionInformation(instrument.Instrument); ok {
					handlePosition(position, instrument.Instrument, enum.D)
					if ok, currentAccountLevel, tradeRule := balanceManager.JudgementProfit(accountManager.GetPosition(instrument.Instrument)); ok {
						distance, trade, position := doCreateNewOrder(instrument, strconv.FormatFloat(accountManager.GetPosition(instrument.Instrument).Units, 'f', 0, 64))
						balanceManager.RegisterNextBalanceManagements(instrument, trade, position, tradeRule, currentAccountLevel, distance)
					}
				}
			}
		case <-tickPerDay.C:
			log.Println("Start Judgement Setup CaptainAmerica SwingTrade :", instrument.Instrument)
			candle := candlesInteractor.GetCandle(dto.NewCandlesGetDto(instrument.Instrument, 2, enum.D))[0]
			candles := doSetupCaptainAmerica(candle, instrument, enum.D)
			tradeStatus, ok := avengers.IsExistSetUpTradeRule(enum.CaptainAmerica, instrument.Instrument, enum.D)
			if ok && captainAmerica.IsExistSecondJudgementTradePlan(instrument.Instrument, enum.D) {
				log.Println("Start Judgement TradePlan CaptainAmerica SwingTrade Second Judge:", instrument.Instrument)
				isOrder, units, captainAmericaStatus := captainAmerica.JudgementTradePlan(tradeStatus, &candles[len(candles)-1], instrument.Instrument, enum.D)
				if isOrder && canCreateNewOrder(instrument.Instrument) {
					// スイングトレードの発注時にデイトレードのポジションを保持している場合はクローズする。
					if ok, position := accountManager.HasPosition(instrument.Instrument); ok {
						accountManager.ClosePosition(instrument.Instrument, position.Units)
					}
					distance, trade, _ := doCreateNewOrder(instrument, units)
					captainAmerica.CompleteTradeStatus(&captainAmericaStatus)
					balanceManager.RegisterFirstBalanceManagements(instrument, trade, enum.CaptainAmerica, distance)
				}
			}
		case <-tickPerHalfDay.C:
			candle := candlesInteractor.GetCandle(dto.NewCandlesGetDto(instrument.Instrument, 2, enum.H12))[0]
			const TimeFormat = "15:04:05"
			if candle.Candles.Time.Format(TimeFormat) == "00:00:00" {
				log.Println("Start Judgement TradePlan CaptainAmerica SwingTrade First Judge:", instrument.Instrument)
				// セットアップの検証時には足データの足種を判定しているため
				tradeStatus, ok := avengers.IsExistSetUpTradeRule(enum.CaptainAmerica, instrument.Instrument, enum.D)
				if ok {
					avengers.JudgementLine(&candle)
					isOrder, units, captainAmericaStatus := captainAmerica.JudgementTradePlan(tradeStatus, &candle, instrument.Instrument, enum.D)
					if isOrder && canCreateNewOrder(instrument.Instrument) {
						// スイングトレードの発注時にデイトレードのポジションを保持している場合はクローズする。
						if ok, position := accountManager.HasPosition(instrument.Instrument); ok {
							accountManager.ClosePosition(instrument.Instrument, position.Units)
						}
						distance, trade, _ := doCreateNewOrder(instrument, units)
						captainAmerica.CompleteTradeStatus(&captainAmericaStatus)
						balanceManager.RegisterFirstBalanceManagements(instrument, trade, enum.CaptainAmerica, distance)
					}
				}
			}
		}
	}
}

// デイトレードのアイアンマンを開始します。
func startDayTradingIronMan() {

}

// スイングトレードのアイアンマンを開始します。
func startSwingTradingIronMan() {

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
func judgementInitialCandles(instrument domain.Instruments, granularity enum.Granularity) []domain.BidAskCandles {
	candles := getCandles(instrument.Instrument, granularity)
	for i := range candles {
		avengers.JudgementLine(&candles[i])
		candles[i].Trend = enum.Range
	}
	avengers.JudgementTrend(&candles[0], &candles[1], &candles[2], &candles[3])
	avengers.CreateTrendStatus(domain.NewTrendStatus(instrument.Instrument, granularity, candles[3].Trend, avengers.Increment(enum.Swing)))
	candlesInteractor.CreateBulkCandles(candles)
	setCandles(candles, instrument.Instrument, granularity)
	return candles
}

// 7つ目以降のポジションとなる場合は、セットアップステータスをリセットします。
func canCreateNewOrder(instrument string) bool {
	account := accountManager.GetAccount()
	if account.OpenPositionCount <= 6 {
		return true
	} else {
		captainAmerica.ResetCaptainAmericaStatus(instrument, enum.H1)
		return false
	}
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

// ポジションの状態と、売買ルールのタイプに応じた処理を行います。
// 共通　　　　；ポジションが0の場合はキャプテンアメリカのセットアップ・取引ステータスをリセットします。
// デイトレード：未実現利益が500円を超えた場合は決済を行います。
func handlePosition(position domain.Positions, instrument string, granularity enum.Granularity) {
	if position.Units == 0 {
		captainAmerica.ResetCaptainAmericaStatus(instrument, granularity)
		log.Print("Reset CaptainAmericaStatus :", instrument, granularity)
	}
	if granularity == enum.H1 && position.UnrealizedPL > config.GetInstance().Property.ProfitGainPrice {
		accountManager.ClosePosition(instrument, position.Units)
	}
}

// キャプテンアメリカのセットアップを検証します。
func doSetupCaptainAmerica(candle domain.BidAskCandles, instrument domain.Instruments, granularity enum.Granularity) []domain.BidAskCandles {
	candles := judgementCurrentCandle(&candle, &instrument, granularity)
	// セットアップの検証
	captainAmerica.JudgementSetup(&candles[len(candles)-2], &candles[len(candles)-1], instrument.Instrument, granularity)
	return candles
}

// トレード計画に従い一連の注文処理を実行します。
// ①キャプテンアメリカステータスを取引中に更新
// ②ポジション情報を生成または更新
// 資金管理レコードの生成に必要な情報を返却します。
func doCreateNewOrder(instrument domain.Instruments, units string) (float64, *domain.Trades, domain.Positions) {
	distance := balanceManager.GetTradingDistance(instrument)
	trade := orderManager.DoNewMarketOrder(instrument.Instrument, units, strconv.FormatFloat(distance, 'f', 5, 64))
	position := accountManager.CreateOrUpdatePosition(instrument.Instrument)
	log.Print("Complete New Order :", instrument.Instrument)
	return distance, trade, position
}
