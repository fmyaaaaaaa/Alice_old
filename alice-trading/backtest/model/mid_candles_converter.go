package model

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/backtest/msg"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"log"
	"strconv"
	"time"
)

type MidCandlesConverter struct{}

// MidCandleのパラメータを想定したSetupRequestをMidCandleに変換します。
func ConvertToMidCandleForSetup(req msg.SetupRequest) MidCandle {
	o := parseFloat(req.OpenPrice)
	c := parseFloat(req.ClosePrice)
	h := parseFloat(req.HighPrice)
	l := parseFloat(req.LowPrice)

	var line enum.Line
	switch {
	case o <= c:
		line = enum.Positive
	case o > c:
		line = enum.Negative
	}

	return MidCandle{
		InstrumentName: req.Instrument,
		Granularity:    enum.Granularity(req.Granularity),
		Open:           o,
		Close:          c,
		High:           h,
		Low:            l,
		Line:           line,
		Time:           stringToTime(req.Time),
	}
}

// MidCandleのパラメータを想定したTradePlanRequestをMidCandleに変換します。
func ConvertToMidCandleForTradePlan(req msg.TradePlanRequest) MidCandle {
	o := parseFloat(req.OpenPrice)
	c := parseFloat(req.ClosePrice)
	h := parseFloat(req.HighPrice)
	l := parseFloat(req.LowPrice)

	var line enum.Line
	switch {
	case o <= c:
		line = enum.Positive
	case o > c:
		line = enum.Negative
	}
	return MidCandle{
		InstrumentName: req.Instrument,
		Granularity:    enum.Granularity(req.Granularity),
		Open:           o,
		Close:          c,
		High:           h,
		Low:            l,
		Line:           line,
		Time:           stringToTime(req.Time),
	}
}

// MidCandleからBidAskCandleに変換します。
func ConvertToBidAskCandle(candle MidCandle) domain.BidAskCandles {
	return domain.BidAskCandles{
		InstrumentName: candle.InstrumentName,
		Granularity:    candle.Granularity,
		Bid: domain.BidRate{
			Open:  candle.Open,
			Close: candle.Close,
			High:  candle.High,
			Low:   candle.Low,
		},
		Ask: domain.AskRate{
			Open:  candle.Open,
			Close: candle.Close,
			High:  candle.High,
			Low:   candle.Low,
		},
		Candles: domain.Candles{
			Time:   candle.Time,
			Volume: 1,
		},
		Line:  candle.Line,
		Trend: candle.Trend,
	}
}

// 文字列をfloat64に変換します。
func parseFloat(str string) float64 {
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Print("fail to parse float64")
	}
	return f
}

func stringToTime(str string) time.Time {
	var layout = "2006-01-02 15:04:05"
	t, _ := time.Parse(layout, str)
	return t
}
