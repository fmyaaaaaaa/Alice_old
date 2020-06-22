package usecase

import (
	"context"
	"fmt"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/api/msg"
	"github.com/fmyaaaaaaa/Alice/alice-trading/usecase/dto"
	"github.com/fmyaaaaaaa/Alice/alice-trading/usecase/util"
	"log"
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
func (c *CandlesInteractor) InitializeCandle(instrumentList []domain.Instruments) {
	for _, instrument := range instrumentList {
		candlesGetDto := dto.NewCandlesGetDto(instrument.Name, 1, enum.H1)
		candleList := c.GetCandle(candlesGetDto)
		c.Create(&candleList[0])
	}
}

// 足データを取得します。
func (c *CandlesInteractor) GetCandle(candlesGetDto dto.CandlesGetDto) []domain.BidAskCandles {
	res, _ := c.Api.GetCandleBidAsk(context.Background(), candlesGetDto.InstrumentName, candlesGetDto.Count, candlesGetDto.Granularity)
	return convertToEntity(res, candlesGetDto.InstrumentName, candlesGetDto.Granularity)
}

// 引数の足データを保存します。
func (c *CandlesInteractor) Create(candle *domain.BidAskCandles) {
	db := c.DB.Connect()
	c.Candles.Create(db, candle)
}

// 主キーをもとに足データを取得します。
func (c *CandlesInteractor) Get(id int) (domain.BidAskCandles, error) {
	db := c.DB.Connect()
	result, err := c.Candles.FindByID(db, id)
	if err != nil {
		log.Print("something wrong")
	}
	return result, nil
}

// 足データを全件取得します。
func (c *CandlesInteractor) GetAll() (candleList []domain.BidAskCandles) {
	db := c.DB.Connect()
	result := c.Candles.FindAll(db)
	return result
}

// 足データを削除します。
// 主キーが必要です。
func (c *CandlesInteractor) Delete(candle *domain.BidAskCandles) {
	db := c.DB.Connect()
	c.Candles.Delete(db, candle)
}

// TODO:LineとTrendを判定するロジックの実装は売買ルールにて対応する。(PositiveとUpTrendの固定値を入れている)
// APIのResponseをBusinessLogicのEntityに変換します。
func convertToEntity(res *msg.CandlesBidAskResponse, instrumentName string, granularity enum.Granularity) []domain.BidAskCandles {
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
				Time:   candle.Time,
				Volume: util.ParseFloat(candle.Volume.String()),
			},
			Line:  enum.Positive,
			Trend: enum.UpTrend,
		}
		fmt.Println(target.InstrumentName, target.Candles.Time)
		result = append(result, *target)
	}
	return result
}
