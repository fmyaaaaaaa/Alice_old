package usecase

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/candles"
	"github.com/fmyaaaaaaa/Alice/alice-trading/infrastructure/config"
	database2 "github.com/fmyaaaaaaa/Alice/alice-trading/infrastructure/database"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/database"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var DB *database2.DB
var candlesInteractor *CandlesInteractor

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	os.Exit(code)
}

func setUp() {
	dummyConf := []string{"", "./../infrastructure/config/env/"}
	config.InitInstance("test", dummyConf)

	DB = database2.NewDB()
	candlesInteractor = &CandlesInteractor{
		DB:      DB,
		Candles: &database.CandlesRepository{},
		Api:     nil,
	}
}

func TestCandlesInteractor_Create_GetAll_Delete(t *testing.T) {
	candle := candles.BidAskCandles{
		InstrumentName: "USD_JPY",
		Granularity:    "H1",
		Bid: candles.BidRate{
			Open:  100,
			Close: 101,
			High:  102,
			Low:   103,
		},
		Ask: candles.AskRate{
			Open:  100,
			Close: 101,
			High:  102,
			Low:   103,
		},
		Candles: candles.Candles{
			Time:   time.Now(),
			Volume: 10,
		},
		Line:  "POSITIVE",
		Trend: "UPDATE",
	}
	candlesInteractor.Create(&candle)
	target := candlesInteractor.GetAll()
	assert.Equal(t, "USD_JPY", target[0].InstrumentName)
	assert.Equal(t, 1, len(target))

	candlesInteractor.Delete(&target[0])
	target = candlesInteractor.GetAll()
	assert.Equal(t, 0, len(target))
}
