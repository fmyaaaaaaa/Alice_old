package rule

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"github.com/fmyaaaaaaa/Alice/alice-trading/usecase"
)

// 売買ルールに共通の処理を管理
type Avengers struct {
	DB              usecase.DBRepository
	Candles         usecase.CandlesRepository
	TradeRuleStatus usecase.TradeRuleStatusRepository
}

// 陽線/陰線を判定します。
func (a Avengers) JudgementLine(currentCandle *domain.BidAskCandles, lastCandle *domain.BidAskCandles) {
	openMid := (currentCandle.Bid.Open + currentCandle.Ask.Open) / 2
	closeMid := (currentCandle.Bid.Close + currentCandle.Ask.Close) / 2
	switch {
	case openMid < closeMid:
		currentCandle.Line = enum.Positive
	case openMid == closeMid:
		currentCandle.Line = lastCandle.Line
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

// TradeRuleStatusを取得します。
func (a Avengers) GetTradeRuleStatus(tradeRule enum.TradeRule, instrument string, granularity enum.Granularity) domain.TradeRuleStatus {
	DB := a.DB.Connect()
	return a.TradeRuleStatus.FindTargetByTradeRuleAndInstrumentAndGranularity(DB, tradeRule, instrument, granularity)
}
