package money

import (
	"context"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"github.com/fmyaaaaaaa/Alice/alice-trading/infrastructure/config"
	"github.com/fmyaaaaaaa/Alice/alice-trading/usecase"
	"github.com/fmyaaaaaaa/Alice/alice-trading/usecase/util"
	"log"
	"math"
)

// 資金管理
type BalanceManager struct {
	DB                 usecase.DBRepository
	PricesApi          usecase.PricesApi
	BackTestResult     usecase.BackTestResultRepository
	BalanceManagements usecase.BalanceManagementsRepository
}

// TODO:円転がbidレートを使うこと確認
// トレーリングストップ値幅を取得します。
// 現在レートと売買をもとに算出します。
func (b BalanceManager) GetTradingDistance(instrument domain.Instruments) float64 {
	// FIXME:lotの単位に応じて基軸通貨レートの値を対応させる
	// 基軸通貨レート
	var distance float64
	// リスク許容額
	riskTolerancePrice := config.GetInstance().Property.RiskTolerancePrice
	if instrument.IsJpyEvaluation() {
		basePrice := 100.0
		distance = b.CalculationPips(float64(riskTolerancePrice), basePrice) / 100
	} else {
		bid, _ := b.getCurrencyPrice(instrument.EvaluationInstrument)
		distance = math.Trunc(b.CalculationPips(float64(riskTolerancePrice), bid)) / 10000
	}
	return distance
}

// 注文結果から初回の資金管理レコードを登録します。
// 初回のため、現在口座水準は0で登録します。
func (b BalanceManager) RegisterFirstBalanceManagements(instrument domain.Instruments, trade *domain.Trades, tradeRule enum.TradeRule, distance float64) {
	nextAccountLevel, delta := b.createNextAccountLevel(tradeRule, instrument, 1, 0)
	balanceManagement := domain.NewBalanceManagements(trade.TradeID, instrument.Instrument, tradeRule, 0, nextAccountLevel, trade.Units, trade.Price, distance, delta)
	b.CreateBalanceManagement(balanceManagement)
}

// 口座残高の定期更新内容から固定比率を検証します。
// ポジションの未実現損益をもとに、次回口座水準の値と比較します。
func (b BalanceManager) JudgementProfit(position domain.Positions) (bool, float64, enum.TradeRule) {
	balanceManagement := b.GetBalanceManagement(position.Instrument)
	nextAccountLevel := balanceManagement.NextAccountLevel
	if nextAccountLevel <= position.UnrealizedPL {
		log.Println("Add Position : ", position.Instrument)
		return true, nextAccountLevel, balanceManagement.TradeRule
	}
	return false, nextAccountLevel, ""
}

// 2回目以降の資金管理レコードを登録します。
func (b BalanceManager) RegisterNextBalanceManagements(instrument domain.Instruments, trade *domain.Trades, position domain.Positions, tradeRule enum.TradeRule, currentAccountLevel, distance float64) {
	buySellUnit := int(math.Abs(position.Units)) / config.GetInstance().Property.OrderLot
	nextAccountLevel, delta := b.createNextAccountLevel(tradeRule, instrument, buySellUnit, currentAccountLevel)
	balanceManagement := domain.NewBalanceManagements(trade.TradeID, instrument.Instrument, tradeRule, currentAccountLevel, nextAccountLevel, trade.Units, trade.Price, distance, delta)
	b.CreateBalanceManagement(balanceManagement)
}

// 損切り値幅を算出します。
// (riskTolerancePrice - (basePrice × price)) / basePrice
func (b BalanceManager) CalculationPips(riskTolerancePrice, basePrice float64) float64 {
	return math.Abs(riskTolerancePrice / basePrice)
}

// 通貨ペアのレートを取得します。
func (b BalanceManager) getCurrencyPrice(instrument string) (bid, ask float64) {
	res := b.PricesApi.GetPrices(context.Background(), instrument)
	bid = util.ParseFloat(res.Prices[0].Bids[0].Price)
	ask = util.ParseFloat(res.Prices[0].Asks[0].Price)
	return bid, ask
}

// 次回口座水準とデルタを算出します。
func (b BalanceManager) createNextAccountLevel(tradeRule enum.TradeRule, instrument domain.Instruments, buySellUnit int, currentAccountLevel float64) (nextAccountLevel, delta float64) {
	backTestResult := b.GetBackTestResult(tradeRule, instrument.Instrument)
	bid, ask := b.getCurrencyPrice(instrument.Instrument)
	mid := bid + ask/2
	baseBid := 0.0
	// デルタ
	if !instrument.IsJpyEvaluation() {
		baseBid, _ = b.getCurrencyPrice(instrument.EvaluationInstrument)
	}
	delta = b.CalculationDelta(instrument, baseBid, backTestResult.MaxDrawDown, mid, instrument.MarginRate)
	// 次回口座水準
	nextAccountLevel = currentAccountLevel + (float64(buySellUnit) * delta)
	return nextAccountLevel, delta
}

// デルタを算出します。
func (b BalanceManager) CalculationDelta(instrument domain.Instruments, basePrice, maxDrawDown, mid, marginRate float64) float64 {
	var minMargin float64
	if instrument.IsJpyEvaluation() {
		minMargin = float64(config.GetInstance().Property.OrderLot) * mid * marginRate
	} else {
		minMargin = float64(config.GetInstance().Property.OrderLot) * mid * marginRate * basePrice
	}
	// FIXME:暫定でこの値にする。
	delta := maxDrawDown + (minMargin / 3)
	return delta
}

// バックテスト結果を取得します。
func (b BalanceManager) GetBackTestResult(tradeRule enum.TradeRule, instrument string) domain.BackTestResults {
	DB := b.DB.Connect()
	return b.BackTestResult.FindByInstrumentAndTradeRule(DB, tradeRule, instrument)
}

// 資金管理を取得します。
// 作成日で降順ソートし、一番目を返却します。
func (b BalanceManager) GetBalanceManagement(instrument string) domain.BalanceManagements {
	DB := b.DB.Connect()
	return b.BalanceManagements.FindByInstrumentOrderByCreatedAt(DB, instrument)
}

// 資金管理を登録します。
func (b BalanceManager) CreateBalanceManagement(balanceManagement *domain.BalanceManagements) {
	DB := b.DB.Connect()
	b.BalanceManagements.Create(DB, balanceManagement)
}
