package rule

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"github.com/fmyaaaaaaa/Alice/alice-trading/usecase"
	"log"
)

type IronMan struct {
	DB                usecase.DBRepository
	TrendStatus       usecase.TrendStatusRepository
	Sequence          usecase.SequenceRepository
	SwingHighLowPrice usecase.SwingHighLowPriceRepository
	SwingTarget       usecase.SwingTargetRepository
	IronManStatus     usecase.IronManStatusRepository
	TradeRuleStatus   usecase.TradeRuleStatusRepository
}

// TODO:呼び出しもとで足データのキャッシュを更新する
// 足データをもとにセットアップを検証します。
func (i IronMan) JudgementSetup(currentCandle *domain.BidAskCandles, instrument string, granularity enum.Granularity) {
	currentTrend := i.GetTrendStatus(instrument, granularity)
	swingID := currentTrend.LastSwingID
	currentHighLowPrice := i.GetHighLowPrice(swingID)

	// トレンドと現在の足データのトレンドが異なる場合、トレンドを更新し、条件に応じてスイングを切り替える。
	// 一緒の場合は高値と安値の比較を行う。
	if currentCandle.Trend != currentTrend.Trend {
		i.createNewSwing(currentCandle, currentTrend, &swingID)
	} else {
		i.handleHighLowPrice(currentCandle, currentHighLowPrice)
	}

	// SwingIDを足データに付与する。
	currentCandle.SwingID = swingID

	// セットアップの検証対象となる高値、安値を取得する。
	swingTarget := i.GetSwingTargetForSetUp(instrument, granularity)
	targetHighLowPrice := i.GetHighLowPrice(swingTarget.SwingID)

	setUp := false
	var ironManStatus *domain.IronManStatus
	switch {
	// 上昇トレンドかつ、仲値がターゲットの高値を超えた場合 または 下降トレンドかつ、仲値がターゲットの安値を超えた場合
	case currentCandle.Trend == enum.UpTrend && targetHighLowPrice.HighPrice <= currentCandle.GetAveMid(),
		currentCandle.Trend == enum.DownTrend && targetHighLowPrice.LowPrice >= currentCandle.GetAveMid():
		setUp = true
	}

	// セットアップ情報を更新する。同一スイングでセットアップを作成済みの場合はスキップする。
	if setUp {
		ironManStatus = domain.NewIronManStatus(instrument, granularity, swingTarget.ID, currentCandle.Trend)
		if ok := i.CreateIronManStatus(ironManStatus); ok {
			tradeRuleStatus := domain.NewTradeRuleStatus(enum.IronMan, instrument, granularity, currentCandle.Candles.Time)
			i.CreateTradeRuleStatus(tradeRuleStatus)
			log.Println("IronMan setup happened ", instrument, granularity)
		}
	}
}

// 足データをもとにトレンドを更新します。
// トレンドの更新時に、下降トレンドに切り替わる場合、スイングの切り替えとして新しいスイングを作成します。
func (i IronMan) createNewSwing(candle *domain.BidAskCandles, currentTrend domain.TrendStatus, swingID *int) {
	params := make(map[string]interface{})
	params["trend"] = candle.Trend
	// 現在のトレンドがレンジ相場かつ、今回足データのトレンドが下降トレンドの場合
	if currentTrend.Trend == enum.Range && candle.Trend == enum.DownTrend {
		nextSwingID := i.increment(enum.Swing)
		*swingID = nextSwingID
		params["last_swing_id"] = *swingID

		// 新しく採番されたSwingIDに対して高値/安値を作成する。
		highLowPrice := domain.NewSwingHighLowPrice(nextSwingID, candle.GetHighMid(), candle.GetLowMid())
		// ターゲットとなる高値/安値を更新する。
		nextSwingTarget := domain.NewSwingTarget(candle.InstrumentName, candle.Granularity, currentTrend.LastSwingID)
		i.CreateHighLowPrice(highLowPrice)
		i.CreateSwingTarget(nextSwingTarget)
	}
	i.UpdateTrendStatus(&currentTrend, params)
}

// 足データの高値/安値がスイングの高値/安値を更新するか検証します。
func (i IronMan) handleHighLowPrice(candle *domain.BidAskCandles, highLowPrice domain.SwingHighLowPrice) {
	// SwingIDの高値/安値に更新があれば行う
	params := make(map[string]interface{})
	switch {
	case highLowPrice.HighPrice < candle.GetHighMid():
		params["high_price"] = candle.GetHighMid()
	case highLowPrice.LowPrice > candle.GetLowMid():
		params["low_price"] = candle.GetLowMid()
	}
	if len(params) != 0 {
		i.UpdateHighLowPrice(&highLowPrice, params)
	}
}

// 足データをもとにトレード計画を判定します。
func (i IronMan) JudgementTradePlan(tradeRuleStatus domain.TradeRuleStatus, candle *domain.BidAskCandles, instrument string, granularity enum.Granularity) {
	// セットアップと同一の足データの場合は処理をスキップ。
	//トレード計画の判定はセットアップの次回足データを対象とするため。
	if tradeRuleStatus.CandleTime.Equal(candle.Candles.Time) {
		return
	}
	// セットアップ情報を取得する。
	ironManStatus := i.GetIronManStatus(instrument, granularity)
	swingTarget := i.GetSwingTargetForTradePlan(ironManStatus.SwingTargetID)
	highLowPrice := i.GetHighLowPrice(swingTarget.SwingID)
	tradePlan := false
	// TODO:資金管理から数量、トレーリングストップ値幅を取得し、OrderManagerから注文を実行する
	switch {
	case ironManStatus.Trend == enum.UpTrend && highLowPrice.HighPrice <= candle.GetAveMid():
		tradePlan = true
		log.Println("IronMan trade happened", instrument, granularity, candle.GetAveMid())
	case ironManStatus.Trend == enum.DownTrend && highLowPrice.LowPrice >= candle.GetAveMid():
		tradePlan = true
		log.Println("IronMan trade happened", instrument, granularity, candle.GetAveMid())
	}
	// セットアップ済みの売買ルールを完了状態にする。
	if tradePlan {
		i.CompleteIronManStatus(&ironManStatus)
		i.CompleteTradeRuleStatus(&tradeRuleStatus)
	}
}

// TrendStatusを作成します。
func (i IronMan) CreateTrendStatus(trendStatus *domain.TrendStatus) {
	DB := i.DB.Connect()
	i.TrendStatus.Create(DB, trendStatus)
}

// TrendStatusを取得します。
func (i IronMan) GetTrendStatus(instrument string, granularity enum.Granularity) domain.TrendStatus {
	DB := i.DB.Connect()
	return i.TrendStatus.FindByInstrumentAndGranularity(DB, instrument, granularity)
}

// TrendStatusを更新します。
func (i IronMan) UpdateTrendStatus(trendStatus *domain.TrendStatus, params map[string]interface{}) {
	DB := i.DB.Connect()
	i.TrendStatus.Update(DB, trendStatus, params)
}

// Swingのシーケンスをインクリメントし、インクリメント後の値を取得します。
func (i IronMan) increment(event enum.Event) int {
	DB := i.DB.Connect()
	return i.Sequence.Increment(DB, event)
}

// SwingHighLowPriceを取得します。
func (i IronMan) GetHighLowPrice(swingID int) domain.SwingHighLowPrice {
	DB := i.DB.Connect()
	return i.SwingHighLowPrice.FindBySwingID(DB, swingID)
}

// SwingHighLowPriceを作成します。
func (i IronMan) CreateHighLowPrice(highLowPrice *domain.SwingHighLowPrice) {
	DB := i.DB.Connect()
	i.SwingHighLowPrice.Create(DB, highLowPrice)
}

// SwingHighLowPriceを更新します。
func (i IronMan) UpdateHighLowPrice(highLowPrice *domain.SwingHighLowPrice, params map[string]interface{}) {
	DB := i.DB.Connect()
	i.SwingHighLowPrice.Update(DB, highLowPrice, params)
}

// SwingTargetを取得します。(セットアップ検証）
func (i IronMan) GetSwingTargetForSetUp(instrument string, granularity enum.Granularity) domain.SwingTarget {
	DB := i.DB.Connect()
	return i.SwingTarget.FindByInstrumentAndGranularity(DB, instrument, granularity)
}

// SwingTargetを取得します。（トレード計画検証）
func (i IronMan) GetSwingTargetForTradePlan(id int) domain.SwingTarget {
	DB := i.DB.Connect()
	return i.SwingTarget.FindByID(DB, id)
}

// SwingTargetを作成します。
func (i IronMan) CreateSwingTarget(swingTarget *domain.SwingTarget) {
	DB := i.DB.Connect()
	i.SwingTarget.Create(DB, swingTarget)
}

// IronManStatusを取得します。
func (i IronMan) GetIronManStatus(instrument string, granularity enum.Granularity) domain.IronManStatus {
	DB := i.DB.Connect()
	return i.IronManStatus.FindByInstrumentAndGranularity(DB, instrument, granularity)
}

// IronManStatusを完了にします。
func (i IronMan) CompleteIronManStatus(ironManStatus *domain.IronManStatus) {
	DB := i.DB.Connect()
	params := map[string]interface{}{
		"status": false,
	}
	i.IronManStatus.Update(DB, ironManStatus, params)
}

// IronManStatusを作成します。
// 同じSwingTargetIDかつ、同一トレンドで既にセットアップ済みの場合はスキップします。
func (i IronMan) CreateIronManStatus(ironManStatus *domain.IronManStatus) bool {
	DB := i.DB.Connect()
	check := i.IronManStatus.FindByInstrumentAndGranularity(DB, ironManStatus.Instrument, ironManStatus.Granularity)
	if check.SwingTargetID == ironManStatus.SwingTargetID && check.Trend == ironManStatus.Trend {
		return false
	}
	i.IronManStatus.Create(DB, ironManStatus)
	return true
}

// TradeRuleStatusを作成します。
func (i IronMan) CreateTradeRuleStatus(tradeRuleStatus *domain.TradeRuleStatus) {
	DB := i.DB.Connect()
	i.TradeRuleStatus.Create(DB, tradeRuleStatus)
}

// TradeRuleStatusを完了として無効にします。
func (i IronMan) CompleteTradeRuleStatus(tradeRuleStatus *domain.TradeRuleStatus) {
	DB := i.DB.Connect()
	params := map[string]interface{}{
		"status": false,
	}
	i.TradeRuleStatus.Update(DB, tradeRuleStatus, params)
}
