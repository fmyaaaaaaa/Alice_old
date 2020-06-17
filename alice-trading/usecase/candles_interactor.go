package usecase

import (
	"context"
	"fmt"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/candles"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/instruments"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/api/msg"
	"github.com/fmyaaaaaaa/Alice/alice-trading/usecase/dto"
	"log"
	"strconv"
)

// 足データのユースケース
type CandlesInteractor struct {
	DB      DBRepository
	Candles CandlesRepository
	Api     CandlesApi
}

// TODO:起動時に取得対象となる足種を指定する。（H1のみ取得している）
// APIで足データを取得し、DBに保存します。
// システム起動時のみ利用します。
func (c *CandlesInteractor) InitializeCandle(instrumentList []instruments.Instruments) {
	for _, instrument := range instrumentList {
		candlesGetDto := dto.NewCandlesGetDto(instrument.Name, 1, enum.H1)
		candleList := c.GetCandle(candlesGetDto)
		c.Create(&candleList[0])
	}
}

// 足データを取得します。
func (c *CandlesInteractor) GetCandle(candlesGetDto dto.CandlesGetDto) []candles.BidAskCandles {
	res, _ := c.Api.GetCandleBidAsk(context.Background(), candlesGetDto.InstrumentName, candlesGetDto.Count, candlesGetDto.Granularity)
	return convertToEntity(res, candlesGetDto.InstrumentName, candlesGetDto.Granularity)
}

// 引数の足データを保存します。
func (c *CandlesInteractor) Create(candle *candles.BidAskCandles) {
	db := c.DB.Connect()
	c.Candles.Create(db, candle)
}

// 主キーをもとに足データを取得します。
func (c *CandlesInteractor) Get(id int) (candles.BidAskCandles, error) {
	db := c.DB.Connect()
	result, err := c.Candles.FindByID(db, id)
	if err != nil {
		log.Print("something wrong")
	}
	return result, nil
}

// 足データを全件取得します。
func (c *CandlesInteractor) GetAll() (candleList []candles.BidAskCandles) {
	db := c.DB.Connect()
	result := c.Candles.FindAll(db)
	return result
}

// 足データを削除します。
// 主キーが必要です。
func (c *CandlesInteractor) Delete(candle *candles.BidAskCandles) {
	db := c.DB.Connect()
	c.Candles.Delete(db, candle)
}

// TODO:LineとTrendを判定するロジックの実装は売買ルールにて対応する。(PositiveとUpTrendの固定値を入れている)
// APIのResponseをBusinessLogicのEntityに変換します。
func convertToEntity(res *msg.CandlesBidAskResponse, instrumentName string, granularity enum.Granularity) []candles.BidAskCandles {
	var result []candles.BidAskCandles
	for _, candle := range res.Candles {
		target := &candles.BidAskCandles{
			InstrumentName: instrumentName,
			Granularity:    granularity,
			Bid: candles.BidRate{
				Open:  parseFloat(candle.Bid.O),
				Close: parseFloat(candle.Bid.C),
				High:  parseFloat(candle.Bid.H),
				Low:   parseFloat(candle.Bid.L),
			},
			Ask: candles.AskRate{
				Open:  parseFloat(candle.Ask.O),
				Close: parseFloat(candle.Ask.C),
				High:  parseFloat(candle.Ask.H),
				Low:   parseFloat(candle.Ask.L),
			},
			Candles: candles.Candles{
				Time:   candle.Time,
				Volume: parseFloat(candle.Volume.String()),
			},
			Line:  enum.POSITIVE,
			Trend: enum.UPTREND,
		}
		fmt.Println(target.InstrumentName, target.Candles.Time)
		result = append(result, *target)
	}
	return result
}

func parseFloat(str string) float64 {
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		panic(err)
	}
	return f
}
