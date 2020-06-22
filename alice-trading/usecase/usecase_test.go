package usecase

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"github.com/fmyaaaaaaa/Alice/alice-trading/infrastructure/config"
	database2 "github.com/fmyaaaaaaa/Alice/alice-trading/infrastructure/database"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/api/oanda"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/database"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var DB *database2.DB
var candlesInteractor *CandlesInteractor
var orderManager *OrderManager

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
	orderManager = &OrderManager{
		DB:              DB,
		Orders:          database.OrdersRepository{},
		Trades:          database.TradesRepository{},
		OrderTradeBinds: database.OrderTradeBindsRepository{},
		OrdersApi:       oanda.NewOrdersApi(),
	}
}

func createDummyOrder() *domain.Orders {
	order := domain.Orders{
		OrderID:     100,
		Instrument:  "USD_JPY",
		Units:       100,
		Type:        "MARKET",
		Price:       100,
		Distance:    1,
		Time:        time.Now(),
		Commission:  0,
		TimeInForce: enum.Fok,
	}
	return &order
}

func createDummyTrade() *domain.Trades {
	trade := domain.Trades{
		TradeID:      101,
		Units:        100,
		Price:        100,
		Instrument:   "USD_JPY",
		State:        "PENDING",
		InitialUnits: 100,
		CurrentUnits: 100,
		RealizedPl:   500,
		UnrealizedPl: 1000,
		MarginUsed:   0,
		OpenTime:     time.Now(),
		CloseTime:    time.Now(),
	}
	return &trade
}

func createDummyBind() *domain.OrderTradeBinds {
	bind := domain.OrderTradeBinds{
		EntryOrderID:    100,
		TradeID:         101,
		StopLossOrderID: 102,
		IsDelete:        false,
	}
	return &bind
}

func TestCandlesInteractor_Create_GetAll_Delete(t *testing.T) {
	candle := domain.BidAskCandles{
		InstrumentName: "USD_JPY",
		Granularity:    "H1",
		Bid: domain.BidRate{
			Open:  100,
			Close: 101,
			High:  102,
			Low:   103,
		},
		Ask: domain.AskRate{
			Open:  100,
			Close: 101,
			High:  102,
			Low:   103,
		},
		Candles: domain.Candles{
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
