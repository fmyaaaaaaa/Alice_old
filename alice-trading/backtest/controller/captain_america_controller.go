package controller

import (
	"encoding/json"
	"github.com/fmyaaaaaaa/Alice/alice-trading/backtest/model"
	"github.com/fmyaaaaaaa/Alice/alice-trading/backtest/msg"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	database2 "github.com/fmyaaaaaaa/Alice/alice-trading/infrastructure/database"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/api/oanda"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/database"
	"github.com/fmyaaaaaaa/Alice/alice-trading/usecase"
	"github.com/fmyaaaaaaa/Alice/alice-trading/usecase/money"
	"github.com/fmyaaaaaaa/Alice/alice-trading/usecase/rule"
	"net/http"
)

var instrumentsInteractor usecase.InstrumentsInteractor
var avengers rule.Avengers
var captainAmerica rule.CaptainAmerica
var balanceManager money.BalanceManager

type CaptainAmericaController struct {
	instrument domain.Instruments
	candles    []domain.BidAskCandles
}

func (c *CaptainAmericaController) StartCaptainAmerica() {
	c.Initialize()
}

// BackTestControllerの実行に必要なインタフェースを初期化します。
func (c *CaptainAmericaController) Initialize() {
	DBRepository := database.DBRepository{DB: database2.NewDB()}
	instrumentsInteractor = usecase.InstrumentsInteractor{
		DB:          &DBRepository,
		Instruments: &database.InstrumentsRepository{},
	}
	avengers = rule.Avengers{
		DB:                &DBRepository,
		TrendStatus:       &database.TrendStatusRepository{},
		Sequence:          &database.SequenceRepository{},
		TradeRuleStatus:   &database.TradeRuleStatusRepository{},
		SwingHighLowPrice: &database.SwingHighLowPriceRepository{},
		SwingTarget:       &database.SwingTargetRepository{},
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
}

// セットアップを検証します。
func (c *CaptainAmericaController) handleSetup(w http.ResponseWriter, r *http.Request) {
	// json parse
	var req msg.SetupRequest
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&req)

	// 足データの判定（線種/トレンド）
	midCandle := model.ConvertToMidCandleForSetup(req)
	candle := model.ConvertToBidAskCandle(midCandle)
	c.addCandle(candle)
	if len(c.candles) > 1 {
		captainAmerica.JudgementSetup(&c.candles[len(c.candles)-2], &c.candles[len(c.candles)-1], candle.InstrumentName, candle.Granularity)
	}

	setup := msg.NewSetupResponse(200, "ok")
	res, err := json.Marshal(setup)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

// トレード計画を検証します。
func (c *CaptainAmericaController) handleTradePlan(w http.ResponseWriter, r *http.Request) {
	// json parse
	var req msg.TradePlanRequest
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&req)

	midCandle := model.ConvertToMidCandleForTradePlan(req)
	candle := model.ConvertToBidAskCandle(midCandle)
	instrument := candle.InstrumentName
	granularity := candle.Granularity

	tradeStatus, ok := avengers.IsExistSetUpTradeRule(enum.CaptainAmerica, candle.InstrumentName, candle.Granularity)
	tradePlan := msg.TradePlanResponse{Status: 200, Result: "OK", IsOrder: false}
	if ok {
		isOrder, units, captainAmericaStatus := captainAmerica.JudgementTradePlan(tradeStatus, &candle, instrument, granularity)
		distance := balanceManager.GetTradingDistance(c.instrument)
		tradePlan = msg.NewTradePlanResponse(200, "ok", units, distance, isOrder)
		captainAmerica.CompleteTradeStatus(&captainAmericaStatus)
		captainAmerica.ResetCaptainAmericaStatus(instrument, granularity)
	}
	res, err := json.Marshal(tradePlan)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

// BackTestControllerが保持するメモリへcandleを追加します。
func (c *CaptainAmericaController) addCandle(candle domain.BidAskCandles) {
	if len(c.candles) == 0 {
		// 初期ロード時に銘柄情報もメモリへ追加する。
		instrument := instrumentsInteractor.GetInstrument(candle.InstrumentName)
		c.instrument = instrument
		// 足データのメモリ追加。
		var candles []domain.BidAskCandles
		candles = append(candles, candle)
		c.candles = candles
	} else {
		c.candles = append(c.candles, candle)
	}
}
