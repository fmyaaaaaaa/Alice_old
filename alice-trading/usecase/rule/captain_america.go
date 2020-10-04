package rule

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"github.com/fmyaaaaaaa/Alice/alice-trading/infrastructure/config"
	"github.com/fmyaaaaaaa/Alice/alice-trading/usecase"
	"log"
	"strconv"
	"time"
)

type CaptainAmerica struct {
	DB                   usecase.DBRepository
	TrendStatus          usecase.TrendStatusRepository
	CaptainAmericaStatus usecase.CaptainAmericaStatusRepository
	TradeRuleStatus      usecase.TradeRuleStatusRepository
}

// 足データをもとにセットアップを検証します。
func (c CaptainAmerica) JudgementSetup(lastCandle, currentCandle *domain.BidAskCandles, instrument string, granularity enum.Granularity) {
	// 同一銘柄、時間足で既に売買中またはセットアップ中の場合はセットアップ検証をスキップする。
	if ok := c.isExistSetupOrTrade(instrument, granularity); ok {
		return
	}
	setUp := false
	// 一つ前の線種と今回の線種を比較し、同一の場合はセットアップ
	if lastCandle.Line == currentCandle.Line {
		setUp = true
	}
	// セットアップ情報を更新する。
	if setUp {
		captainAmericaStatus := domain.NewCaptainAmericaStatus(instrument, granularity, currentCandle.Line, currentCandle.GetCloseMid(), true, false)
		c.CreateOrUpdateCaptainAmericaStatus(captainAmericaStatus)
		tradeRuleStatus := domain.NewTradeRuleStatus(enum.CaptainAmerica, instrument, granularity, currentCandle.Candles.Time.Add(12*time.Hour))
		c.CreateOrUpdateTradeRuleStatus(tradeRuleStatus)
		log.Println("CaptainAmerica setup happened ", instrument, granularity)
	}
}

// 足データをもとにスイングトレードのトレード計画を判定します。
func (c CaptainAmerica) JudgementTradePlanOfSwingTrade(tradeRuleStatus domain.TradeRuleStatus, currentCandle *domain.BidAskCandles, instrument string, granularity enum.Granularity, judgementCandle *domain.BidAskCandles) (bool, string, domain.CaptainAmericaStatus) {
	// 注文数量
	units := 0
	if tradeRuleStatus.CandleTime.Equal(currentCandle.Candles.Time) {
		return false, strconv.Itoa(units), domain.CaptainAmericaStatus{}
	}

	// セットアップを取得
	captainAmericaStatus := c.GetCaptainAmericaStatus(instrument, granularity)
	tradePlan := false
	switch captainAmericaStatus.Line {
	case enum.Positive:
		if captainAmericaStatus.SetupPrice <= currentCandle.GetCloseMid() && judgementCandle.Line == enum.Positive {
			tradePlan = true
			units = config.GetInstance().Property.OrderLot
			log.Println("CaptainAmerica trade happened", currentCandle.Candles.Time, instrument, granularity, currentCandle.GetCloseMid())
		}
	case enum.Negative:
		if captainAmericaStatus.SetupPrice >= currentCandle.GetCloseMid() && judgementCandle.Line == enum.Negative {
			tradePlan = true
			units = -config.GetInstance().Property.OrderLot
			log.Println("CaptainAmerica trade happened", currentCandle.Candles.Time, instrument, granularity, currentCandle.GetCloseMid())
		}
	}
	// トレード計画の結果に応じて、売買ルールの状態を変更する。
	// Second Judgeを実施しないように修正するため、判定が完了したら、セットアップ状況をリセットする。
	//if tradePlan || captainAmericaStatus.SecondJudge {
	c.CompleteTradeRuleStatus(&tradeRuleStatus)
	//}
	c.HandleCaptainAmericaStatus(&captainAmericaStatus, tradePlan)
	return tradePlan, strconv.Itoa(units), captainAmericaStatus
}

// 足データをもとにデイトレードのトレード計画を判定します。
func (c CaptainAmerica) JudgementTradePlanOfDayTrade(tradeRuleStatus domain.TradeRuleStatus, currentCandle, additionalCandle *domain.BidAskCandles, instrument string, granularity enum.Granularity) (bool, string, domain.CaptainAmericaStatus) {
	// 注文数量
	units := 0
	// セットアップと同一の足データの場合は処理をスキップ。
	// トレード計画の判定はセットアップの次回足データを対象とするため。
	if tradeRuleStatus.CandleTime.Equal(currentCandle.Candles.Time) {
		return false, strconv.Itoa(units), domain.CaptainAmericaStatus{}
	}
	// セットアップを取得
	captainAmericaStatus := c.GetCaptainAmericaStatus(instrument, granularity)
	tradePlan := false
	switch captainAmericaStatus.Line {
	case enum.Positive:
		// セットアップ時点の価格をトレード検証の30分足データが上回っている かつ 直前の15分足の線種がPositiveであること。
		if captainAmericaStatus.SetupPrice <= currentCandle.GetCloseMid() && currentCandle.Line == enum.Positive && additionalCandle.Line == enum.Positive {
			tradePlan = true
			units = config.GetInstance().Property.OrderLot
			log.Println("CaptainAmerica trade happened", currentCandle.Candles.Time, instrument, granularity, currentCandle.GetCloseMid())
		}
	case enum.Negative:
		// セットアップ時点の価格をトレード検証の30分足データが下回っている かつ 直前の15分足の線種がNegativeであること。
		if captainAmericaStatus.SetupPrice >= currentCandle.GetCloseMid() && currentCandle.Line == enum.Negative && additionalCandle.Line == enum.Negative {
			tradePlan = true
			units = -config.GetInstance().Property.OrderLot
			log.Println("CaptainAmerica trade happened", currentCandle.Candles.Time, instrument, granularity, currentCandle.GetCloseMid())
		}
	}
	// トレード計画の結果に応じて、売買ルールの状態を変更する。
	if tradePlan || captainAmericaStatus.SecondJudge {
		c.CompleteTradeRuleStatus(&tradeRuleStatus)
	}
	c.HandleCaptainAmericaStatus(&captainAmericaStatus, tradePlan)
	return tradePlan, strconv.Itoa(units), captainAmericaStatus
}

// 銘柄、足種でセットアップ済みまたは取引済みかどうかを確認します。
// セットアップ済み、取引済みの場合はtrueを返却します。
func (c CaptainAmerica) isExistSetupOrTrade(instrument string, granularity enum.Granularity) bool {
	captainAmericaStatus := c.GetCaptainAmericaStatus(instrument, granularity)
	if captainAmericaStatus.SetupStatus || captainAmericaStatus.TradeStatus {
		return true
	}
	return false
}

// 2回目のトレード計画検証対象が存在するかどうかを確認します。
// SecondJudge対象（true）のセットアップが存在する場合はtrueを返却します。
func (c CaptainAmerica) IsExistSecondJudgementTradePlan(instrument string, granularity enum.Granularity) bool {
	captainAmericaStatus := c.GetCaptainAmericaStatus(instrument, granularity)
	return captainAmericaStatus.SecondJudge
}

// CaptainAmericaStatusを取得します。
func (c CaptainAmerica) GetCaptainAmericaStatus(instrument string, granularity enum.Granularity) domain.CaptainAmericaStatus {
	DB := c.DB.Connect()
	return c.CaptainAmericaStatus.FindByInstrumentAndGranularity(DB, instrument, granularity)
}

// CaptainAmericaStatusを作成します。
// 既に同一銘柄、足種で作成済みの場合は更新します。
func (c CaptainAmerica) CreateOrUpdateCaptainAmericaStatus(captainAmericaStatus *domain.CaptainAmericaStatus) {
	DB := c.DB.Connect()
	if target := c.CaptainAmericaStatus.FindByInstrumentAndGranularity(DB, captainAmericaStatus.Instrument, captainAmericaStatus.Granularity); target.ID == 0 {
		c.CaptainAmericaStatus.Create(DB, captainAmericaStatus)
	} else {
		params := map[string]interface{}{
			"line":         captainAmericaStatus.Line,
			"setup_price":  captainAmericaStatus.SetupPrice,
			"setup_status": captainAmericaStatus.SetupStatus,
			"trade_status": captainAmericaStatus.TradeStatus,
		}
		c.CaptainAmericaStatus.Update(DB, &target, params)
	}
}

// キャプテンアメリカのステータスをリセットします。
func (c CaptainAmerica) ResetCaptainAmericaStatus(instrument string, granularity enum.Granularity) {
	DB := c.DB.Connect()
	c.CaptainAmericaStatus.Reset(DB, instrument, granularity)
}

// 新規注文完了後にトレードステータスを更新します。
func (c CaptainAmerica) CompleteTradeStatus(captainAmericaStatus *domain.CaptainAmericaStatus) {
	DB := c.DB.Connect()
	// セットアップ済みの売買ルールを完了、取引ステータスを取引中に更新する。
	params := map[string]interface{}{
		"trade_status": true,
	}
	c.CaptainAmericaStatus.Update(DB, captainAmericaStatus, params)
}

// トレード計画の検証結果に応じて、キャプテンアメリカのステータスを更新します。
func (c CaptainAmerica) HandleCaptainAmericaStatus(captainAmericaStatus *domain.CaptainAmericaStatus, tradePlan bool) {
	DB := c.DB.Connect()
	params := make(map[string]interface{})
	// セットアップ済みの売買ルールを完了、取引ステータスを取引中に更新する。
	//if tradePlan {
	params["setup_status"] = false
	params["second_judge"] = false
	//} else {
	//	if captainAmericaStatus.SecondJudge {
	//		params["second_judge"] = false
	//		params["setup_status"] = false
	//	} else {
	//		params["second_judge"] = false
	//	}
	//}
	c.CaptainAmericaStatus.Update(DB, captainAmericaStatus, params)
}

// TradeRuleStatusを作成します。
// 既に同一銘柄、足種で作成済みの場合は更新します。
func (c CaptainAmerica) CreateOrUpdateTradeRuleStatus(tradeRuleStatus *domain.TradeRuleStatus) {
	DB := c.DB.Connect()
	if target := c.TradeRuleStatus.FindByTradeRuleAndInstrumentAndGranularity(DB, tradeRuleStatus.TradeRule, tradeRuleStatus.Instrument, tradeRuleStatus.Granularity); target.ID == 0 {
		c.TradeRuleStatus.Create(DB, tradeRuleStatus)
	} else {
		params := map[string]interface{}{
			"candle_time": tradeRuleStatus.CandleTime,
			"status":      tradeRuleStatus.Status,
		}
		c.TradeRuleStatus.Update(DB, &target, params)
	}
}

// TradeRuleStatusを完了として無効にします。
func (c CaptainAmerica) CompleteTradeRuleStatus(tradeRuleStatus *domain.TradeRuleStatus) {
	DB := c.DB.Connect()
	params := map[string]interface{}{
		"status": false,
	}
	c.TradeRuleStatus.Update(DB, tradeRuleStatus, params)
}
