package rule

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"github.com/fmyaaaaaaa/Alice/alice-trading/usecase"
)

// 売買ルールに共通の処理を管理
type Avengers struct {
	DB                usecase.DBRepository
	TrendStatus       usecase.TrendStatusRepository
	Sequence          usecase.SequenceRepository
	TradeRuleStatus   usecase.TradeRuleStatusRepository
	SwingHighLowPrice usecase.SwingHighLowPriceRepository
	SwingTarget       usecase.SwingTargetRepository
}

// 陽線/陰線を判定します。
func (a Avengers) JudgementLine(currentCandle *domain.BidAskCandles) {
	openMid := (currentCandle.Bid.Open + currentCandle.Ask.Open) / 2
	closeMid := (currentCandle.Bid.Close + currentCandle.Ask.Close) / 2
	switch {
	case openMid <= closeMid:
		currentCandle.Line = enum.Positive
	case openMid > closeMid:
		currentCandle.Line = enum.Negative
	}
}

// 上昇トレンド/下降トレンド/レンジ相場を判定します。
func (a Avengers) JudgementTrend(threeLineAgoCandle, twoLineAgoCandle, lastCandle, currentCandle *domain.BidAskCandles) {
	var candles []domain.BidAskCandles
	candles = append(candles, *threeLineAgoCandle, *twoLineAgoCandle, *lastCandle, *currentCandle)
	currentCandle.Trend = countCandleLine(candles)
}

// 足データをもとにスイング判定、相場全体のトレンド判定を行います。
func (a Avengers) JudgementSwingAndAllTrend(currentCandle *domain.BidAskCandles, instrument string, granularity enum.Granularity) {
	currentTrend := a.GetTrendStatus(instrument, granularity)
	swingID := currentTrend.LastSwingID
	currentHighLowPrice := a.GetHighLowPrice(swingID)

	// トレンドと現在の足データのトレンドが異なる場合、トレンドを更新し、条件に応じてスイングを切り替える。
	// 一緒の場合は高値と安値の比較を行う。
	if currentCandle.Trend != currentTrend.Trend {
		a.createNewSwing(currentCandle, currentTrend, &swingID)
	} else {
		a.handleHighLowPrice(currentCandle, currentHighLowPrice)
	}
	// SwingIDを足データに付与する。
	currentCandle.SwingID = swingID
}

// 売買ルールに対してセットアップが存在するかどうかを確認します。
// 存在する場合は、第二戻り値でtrueを返却します。
func (a Avengers) IsExistSetUpTradeRule(tradeRule enum.TradeRule, instrument string, granularity enum.Granularity) (domain.TradeRuleStatus, bool) {
	tradeRuleStatus := a.GetTradeRuleStatus(tradeRule, instrument, granularity)
	if tradeRuleStatus.ID == 0 {
		return tradeRuleStatus, false
	}
	return tradeRuleStatus, true
}

// 直近4本の陽線/陰線からトレンドを判定します。
func countCandleLine(candles []domain.BidAskCandles) enum.Trend {
	var positive, negative int
	for _, candle := range candles {
		switch candle.Line {
		case enum.Positive:
			positive++
		case enum.Negative:
			negative++
		}
	}
	switch {
	case positive >= 3:
		return enum.UpTrend
	case negative >= 3:
		return enum.DownTrend
	default:
		return enum.Range
	}
}

// 足データをもとにトレンドを更新します。
// トレンドの更新時に、下降トレンドに切り替わる場合、スイングの切り替えとして新しいスイングを作成します。
func (a Avengers) createNewSwing(candle *domain.BidAskCandles, currentTrend domain.TrendStatus, swingID *int) {
	params := make(map[string]interface{})
	params["trend"] = candle.Trend
	// 現在のトレンドがレンジ相場かつ、今回足データのトレンドが下降トレンドの場合
	if currentTrend.Trend == enum.Range && candle.Trend == enum.DownTrend {
		nextSwingID := a.Increment(enum.Swing)
		*swingID = nextSwingID
		params["last_swing_id"] = *swingID

		// 新しく採番されたSwingIDに対して高値/安値を作成する。
		highLowPrice := domain.NewSwingHighLowPrice(nextSwingID, candle.GetHighMid(), candle.GetLowMid())
		// ターゲットとなる高値/安値を更新する。
		nextSwingTarget := domain.NewSwingTarget(candle.InstrumentName, candle.Granularity, currentTrend.LastSwingID)
		a.CreateHighLowPrice(highLowPrice)
		a.CreateSwingTarget(nextSwingTarget)
	}
	a.UpdateTrendStatus(&currentTrend, params)
}

// 足データの高値/安値がスイングの高値/安値を更新するか検証します。
func (a Avengers) handleHighLowPrice(candle *domain.BidAskCandles, highLowPrice domain.SwingHighLowPrice) {
	// SwingIDの高値/安値に更新があれば行う
	params := make(map[string]interface{})
	switch {
	case highLowPrice.HighPrice < candle.GetHighMid():
		params["high_price"] = candle.GetHighMid()
	case highLowPrice.LowPrice > candle.GetLowMid():
		params["low_price"] = candle.GetLowMid()
	}
	if len(params) != 0 {
		a.UpdateHighLowPrice(&highLowPrice, params)
	}
}

// TrendStatusを作成します。
func (a Avengers) CreateTrendStatus(trendStatus *domain.TrendStatus) {
	DB := a.DB.Connect()
	a.TrendStatus.Create(DB, trendStatus)
}

// TrendStatusを取得します。
func (a Avengers) GetTrendStatus(instrument string, granularity enum.Granularity) domain.TrendStatus {
	DB := a.DB.Connect()
	return a.TrendStatus.FindByInstrumentAndGranularity(DB, instrument, granularity)
}

// TrendStatusを更新します。
func (a Avengers) UpdateTrendStatus(trendStatus *domain.TrendStatus, params map[string]interface{}) {
	DB := a.DB.Connect()
	a.TrendStatus.Update(DB, trendStatus, params)
}

// SwingHighLowPriceを取得します。
func (a Avengers) GetHighLowPrice(swingID int) domain.SwingHighLowPrice {
	DB := a.DB.Connect()
	return a.SwingHighLowPrice.FindBySwingID(DB, swingID)
}

// SwingHighLowPriceを作成します。
func (a Avengers) CreateHighLowPrice(highLowPrice *domain.SwingHighLowPrice) {
	DB := a.DB.Connect()
	a.SwingHighLowPrice.Create(DB, highLowPrice)
}

// SwingHighLowPriceを更新します。
func (a Avengers) UpdateHighLowPrice(highLowPrice *domain.SwingHighLowPrice, params map[string]interface{}) {
	DB := a.DB.Connect()
	a.SwingHighLowPrice.Update(DB, highLowPrice, params)
}

// Swingのシーケンスをインクリメントし、インクリメント後の値を取得します。
func (a Avengers) Increment(event enum.Event) int {
	DB := a.DB.Connect()
	return a.Sequence.Increment(DB, event)
}

// SwingTargetを作成します。
func (a Avengers) CreateSwingTarget(swingTarget *domain.SwingTarget) {
	DB := a.DB.Connect()
	a.SwingTarget.Create(DB, swingTarget)
}

// TradeRuleStatusを取得します。
func (a Avengers) GetTradeRuleStatus(tradeRule enum.TradeRule, instrument string, granularity enum.Granularity) domain.TradeRuleStatus {
	DB := a.DB.Connect()
	return a.TradeRuleStatus.FindTargetByTradeRuleAndInstrumentAndGranularity(DB, tradeRule, instrument, granularity)
}
