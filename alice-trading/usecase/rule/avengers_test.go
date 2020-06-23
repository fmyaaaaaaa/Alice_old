package rule

import (
	"encoding/json"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"github.com/fmyaaaaaaa/Alice/alice-trading/infrastructure/config"
	"github.com/fmyaaaaaaa/Alice/alice-trading/infrastructure/database"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/api/msg"
	database2 "github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/database"
	"github.com/fmyaaaaaaa/Alice/alice-trading/usecase/util"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"
)

var avengers Avengers
var ironMan IronMan

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	os.Exit(code)
}

func setUp() {
	dummyConf := []string{"", "./../../infrastructure/config/env/"}
	config.InitInstance("test", dummyConf)
	DB := database.NewDB()
	avengers = Avengers{
		DB:              &database2.DBRepository{DB: DB},
		Candles:         &database2.CandlesRepository{},
		TradeRuleStatus: &database2.TradeRuleStatusRepository{},
	}
	ironMan = IronMan{
		DB:                &database2.DBRepository{DB: DB},
		TrendStatus:       &database2.TrendStatusRepository{},
		Sequence:          &database2.SequenceRepository{},
		SwingHighLowPrice: &database2.SwingHighLowPriceRepository{},
		SwingTarget:       &database2.SwingTargetRepository{},
		IronManStatus:     &database2.IronManStatusRepository{},
		TradeRuleStatus:   &database2.TradeRuleStatusRepository{},
	}
}

func TestAvengers_JudgementLine(t *testing.T) {
	// 一つ前の足データ
	lastCandle := &domain.BidAskCandles{
		Bid:  domain.BidRate{Open: 100, Close: 120},
		Ask:  domain.AskRate{Open: 110, Close: 130},
		Line: enum.Positive,
	}
	// 陽線
	positiveCandle := &domain.BidAskCandles{
		Bid: domain.BidRate{Open: 110, Close: 120},
		Ask: domain.AskRate{Open: 120, Close: 130},
	}
	// 陰線
	negativeCandle := &domain.BidAskCandles{
		Bid: domain.BidRate{Open: 120, Close: 110},
		Ask: domain.AskRate{Open: 130, Close: 120},
	}
	// 前足データ線種
	sameCandle := &domain.BidAskCandles{
		Bid: domain.BidRate{Open: 110, Close: 110},
		Ask: domain.AskRate{Open: 120, Close: 120},
	}
	avengers.JudgementLine(positiveCandle, lastCandle)
	avengers.JudgementLine(negativeCandle, lastCandle)
	avengers.JudgementLine(sameCandle, lastCandle)

	assert.Equal(t, enum.Positive, positiveCandle.Line)
	assert.Equal(t, enum.Negative, negativeCandle.Line)
	assert.Equal(t, enum.Positive, sameCandle.Line)
}

func TestAvengers_JudgementTrend(t *testing.T) {
	// 元データ
	positiveCandle1 := &domain.BidAskCandles{
		Bid:  domain.BidRate{Open: 100, Close: 120},
		Ask:  domain.AskRate{Open: 110, Close: 130},
		Line: enum.Positive,
	}
	positiveCandle2 := &domain.BidAskCandles{
		Bid:  domain.BidRate{Open: 100, Close: 120},
		Ask:  domain.AskRate{Open: 110, Close: 130},
		Line: enum.Positive,
	}
	positiveCandle3 := &domain.BidAskCandles{
		Bid:  domain.BidRate{Open: 100, Close: 120},
		Ask:  domain.AskRate{Open: 110, Close: 130},
		Line: enum.Positive,
	}
	negativeCandle1 := &domain.BidAskCandles{
		Bid:  domain.BidRate{Open: 120, Close: 110},
		Ask:  domain.AskRate{Open: 130, Close: 120},
		Line: enum.Negative,
	}
	negativeCandle2 := &domain.BidAskCandles{
		Bid:  domain.BidRate{Open: 120, Close: 110},
		Ask:  domain.AskRate{Open: 130, Close: 120},
		Line: enum.Negative,
	}
	negativeCandle3 := &domain.BidAskCandles{
		Bid:  domain.BidRate{Open: 120, Close: 110},
		Ask:  domain.AskRate{Open: 130, Close: 120},
		Line: enum.Negative,
	}

	// テスト対象データ（陽線×４）
	targetPositiveCandle1 := &domain.BidAskCandles{
		Bid:  domain.BidRate{Open: 100, Close: 120},
		Ask:  domain.AskRate{Open: 110, Close: 130},
		Line: enum.Positive,
	}
	avengers.JudgementTrend(positiveCandle1, positiveCandle2, positiveCandle3, targetPositiveCandle1)
	assert.Equal(t, enum.UpTrend, targetPositiveCandle1.Trend)

	// テスト対象データ（陽線×３ 陰線×１）
	targetPositiveCandle2 := &domain.BidAskCandles{
		Bid:  domain.BidRate{Open: 100, Close: 120},
		Ask:  domain.AskRate{Open: 110, Close: 130},
		Line: enum.Positive,
	}
	avengers.JudgementTrend(positiveCandle1, positiveCandle2, negativeCandle1, targetPositiveCandle2)
	assert.Equal(t, enum.UpTrend, targetPositiveCandle2.Trend)

	// テスト対象データ（陽線×２ 陰線×２）
	targetPositiveCandle3 := &domain.BidAskCandles{
		Bid:  domain.BidRate{Open: 100, Close: 120},
		Ask:  domain.AskRate{Open: 110, Close: 130},
		Line: enum.Positive,
	}
	avengers.JudgementTrend(negativeCandle1, positiveCandle1, negativeCandle2, targetPositiveCandle3)
	assert.Equal(t, enum.Range, targetPositiveCandle3.Trend)

	// テスト対象データ（陽線×１ 陰線×３）
	targetPositiveCandle4 := &domain.BidAskCandles{
		Bid:  domain.BidRate{Open: 100, Close: 120},
		Ask:  domain.AskRate{Open: 110, Close: 130},
		Line: enum.Positive,
	}
	avengers.JudgementTrend(negativeCandle1, negativeCandle2, negativeCandle3, targetPositiveCandle4)
	assert.Equal(t, enum.DownTrend, targetPositiveCandle4.Trend)

	// テスト対象データ（陰線×４）
	targetNegativeCandle1 := &domain.BidAskCandles{
		Bid:  domain.BidRate{Open: 120, Close: 110},
		Ask:  domain.AskRate{Open: 130, Close: 120},
		Line: enum.Negative,
	}
	avengers.JudgementTrend(negativeCandle1, negativeCandle2, negativeCandle3, targetNegativeCandle1)
	assert.Equal(t, enum.DownTrend, targetNegativeCandle1.Trend)

	// テスト対象データ（陰線×３ 陽線×１）
	targetNegativeCandle2 := &domain.BidAskCandles{
		Bid:  domain.BidRate{Open: 120, Close: 110},
		Ask:  domain.AskRate{Open: 130, Close: 120},
		Line: enum.Negative,
	}
	avengers.JudgementTrend(negativeCandle1, positiveCandle1, negativeCandle2, targetNegativeCandle2)
	assert.Equal(t, enum.DownTrend, targetNegativeCandle2.Trend)

	// テスト対象データ（陰線×２ 陽線×２）
	targetNegativeCandle3 := &domain.BidAskCandles{
		Bid:  domain.BidRate{Open: 120, Close: 110},
		Ask:  domain.AskRate{Open: 130, Close: 120},
		Line: enum.Negative,
	}
	avengers.JudgementTrend(positiveCandle1, positiveCandle2, negativeCandle1, targetNegativeCandle3)
	assert.Equal(t, enum.Range, targetNegativeCandle3.Trend)

	// テスト対象データ（陰線×１ 陽線×３）
	targetNegativeCandle4 := &domain.BidAskCandles{
		Bid:  domain.BidRate{Open: 120, Close: 110},
		Ask:  domain.AskRate{Open: 130, Close: 120},
		Line: enum.Negative,
	}
	avengers.JudgementTrend(positiveCandle1, positiveCandle2, positiveCandle3, targetNegativeCandle4)
	assert.Equal(t, enum.UpTrend, targetNegativeCandle4.Trend)
}

func TestIronMan_JudgementSetup_TradePlan(t *testing.T) {
	initTestDataForIronMan()
	candles := createTestCandlesAUDNZD()
	for i := range candles {
		ironMan.JudgementSetup(&candles[i], "AUD_NZD", enum.H1)
		tradeRuleStatus, ok := avengers.IsExistSetUpTradeRule(enum.IronMan, "AUD_NZD", enum.H1)
		if ok {
			ironMan.JudgementTradePlan(tradeRuleStatus, &candles[i], "AUD_NZD", enum.H1)
		}
	}
	ironManStatus := ironMan.GetIronManStatus("AUD_NZD", enum.H1)
	assert.Equal(t, false, ironManStatus.Status)
}

// Swingを考慮した足データのテストデータ（AUD_NZD 2020-06-20 00:00:00 ~ 2020-06-23 19:00:00)
func createTestCandlesAUDNZD() []domain.BidAskCandles {
	data, err := os.Open("./../../data/candles_AUD_NZD.json")
	if err != nil {
		log.Print(err)
	}
	defer data.Close()
	jsonDecoder := json.NewDecoder(data)
	var candlesRes msg.CandlesBidAskResponse
	if err = jsonDecoder.Decode(&candlesRes); err != nil {
		log.Println(err)
	}
	candles := convertToEntity(&candlesRes, "AUD_NZD", enum.H1)
	// 最初の４本は手動にて線種、トレンドを設定
	candles[0].Line = enum.Positive
	candles[0].SwingID = 1
	candles[0].Trend = enum.Range
	candles[1].Trend = enum.Range
	candles[2].Trend = enum.UpTrend
	candles[3].Trend = enum.UpTrend
	for i := range candles {
		if i == 0 {
			continue
		} else {
			avengers.JudgementLine(&candles[i], &candles[i-1])
		}
		if i == 1 || i == 2 || i == 3 {
			continue
		} else {
			avengers.JudgementTrend(&candles[i-3], &candles[i-2], &candles[i-1], &candles[i])
		}
		// FIXME:DEBUG用にコメントアウトしておくが、不要になれば削除する。
		//fmt.Println(fmt.Sprintf("CandleTime: %s, Line: %s, Trend: %s, OpenMid: %s, CloseMid: %s, Open-Close: %s",
		//	candles[i].Candles.Time, candles[i].Line, candles[i].Trend,
		//	strconv.FormatFloat(candles[i].GetOpenMid(), 'f', -1, 64),
		//	strconv.FormatFloat(candles[i].GetCloseMid(), 'f', -1, 64),
		//	strconv.FormatFloat(candles[i].GetOpenMid()-candles[i].GetCloseMid(), 'f', -1, 64)))
	}
	return candles
}

// セットアップの検証用の初期データを作成する。
func initTestDataForIronMan() {
	// トレンドステータス
	trendStatus := domain.NewTrendStatus("AUD_NZD", enum.H1, enum.Range, 1)
	// 高値/安値
	highLowPrice := domain.NewSwingHighLowPrice(1, 1.07076, 1.06811)
	// スイングターゲット
	swingTarget := domain.NewSwingTarget("AUD_NZD", enum.H1, 1)

	ironMan.CreateTrendStatus(trendStatus)
	ironMan.CreateHighLowPrice(highLowPrice)
	ironMan.CreateSwingTarget(swingTarget)
}

// 足データ変換処理のDeepCopy
func convertToEntity(res *msg.CandlesBidAskResponse, instrumentName string, granularity enum.Granularity) []domain.BidAskCandles {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	var result []domain.BidAskCandles
	for _, candle := range res.Candles {
		target := &domain.BidAskCandles{
			InstrumentName: instrumentName,
			Granularity:    granularity,
			Bid: domain.BidRate{
				Open:  util.ParseFloat(candle.Bid.O),
				Close: util.ParseFloat(candle.Bid.C),
				High:  util.ParseFloat(candle.Bid.H),
				Low:   util.ParseFloat(candle.Bid.L),
			},
			Ask: domain.AskRate{
				Open:  util.ParseFloat(candle.Ask.O),
				Close: util.ParseFloat(candle.Ask.C),
				High:  util.ParseFloat(candle.Ask.H),
				Low:   util.ParseFloat(candle.Ask.L),
			},
			Candles: domain.Candles{
				Time:   candle.Time.In(jst),
				Volume: util.ParseFloat(candle.Volume.String()),
			},
		}
		result = append(result, *target)
	}
	return result
}
