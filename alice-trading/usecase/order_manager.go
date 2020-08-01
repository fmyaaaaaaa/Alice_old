package usecase

import (
	"context"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"github.com/fmyaaaaaaa/Alice/alice-trading/infrastructure/config"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/api/msg"
	"github.com/fmyaaaaaaa/Alice/alice-trading/usecase/util"
	"log"
	"strconv"
)

// 注文管理
type OrderManager struct {
	DB              DBRepository
	Orders          OrdersRepository
	Trades          TradesRepository
	OrderTradeBinds OrderTradeBindsRepository
	OrdersApi       OrdersApi
	TradesApi       TradesApi
}

// 新規Market注文を発注します。
// 発注時にトレーリングストップを設定します。
// 発注結果から新規注文、トレーリング注文、取引情報を保存します。
func (o OrderManager) DoNewMarketOrderTrailingStop(instrument, units, distance string) *domain.Trades {
	// トレーリングストップ
	trailing := &msg.TrailingStopLossDetails{
		Distance:    distance,
		TimeInForce: enum.Gtc,
	}
	// 注文リクエスト実行
	reqParam := msg.NewMarketOrderRequest(instrument, units, "", enum.DefaultOrderPositionFIll, nil, nil, trailing)
	createRes, errRes := o.OrdersApi.CreateNewOrder(context.Background(), reqParam)
	if errRes != nil {
		log.Print("fail to new Order", errRes.ErrorCode, errRes.ErrorMessage)
		return nil
	}
	// 注文情報を保存する
	o.CreateOrder(o.convertToEntityCreateOrder(createRes, enum.Fok, enum.Market, distance))
	// 取引情報を保存する
	trade := o.convertToEntityTrade(createRes)
	o.CreateTrade(trade)
	// トレーリングストップ注文情報を取得し、保存する。
	// FIXME:注文結果のレスポンスでトレーリングストップのIDに合わせて確認する
	getRes := o.OrdersApi.GetOrder(context.Background(), createRes.LastTransactionID)
	o.CreateOrder(o.convertToEntityGetOrder(getRes, instrument, units))
	// 注文情報と取引情報の紐付けを保存する。
	o.CreateBind(o.convertToEntityBind(createRes))

	return trade
}

// 新規Market注文を発注します。
// 発注時にストップリミットを設定します。
// 発注結果から新規注文、ストップリミット注文、取引情報を保存します。
func (o OrderManager) DoNewMarketOrderStopLimit(instrument, units, distance string) *domain.Trades {
	// ストップリミット
	stopLoss := &msg.StopLossDetails{
		Distance:    distance,
		TimeInForce: enum.Gtc,
	}
	// 注文リクエスト実行
	reqParam := msg.NewMarketOrderRequest(instrument, units, "", enum.DefaultOrderPositionFIll, nil, stopLoss, nil)
	createRes, errRes := o.OrdersApi.CreateNewOrder(context.Background(), reqParam)
	if errRes != nil {
		log.Print("fail to new Order", errRes.ErrorCode, errRes.ErrorMessage)
		return nil
	}
	o.CreateOrder(o.convertToEntityCreateOrder(createRes, enum.Fok, enum.Market, distance))
	trade := o.convertToEntityTrade(createRes)
	getRes := o.OrdersApi.GetOrder(context.Background(), createRes.LastTransactionID)
	o.CreateOrder(o.convertToEntityGetOrderForStopLimit(getRes, instrument, units))
	o.CreateBind(o.convertToEntityBind(createRes))
	return trade
}

// 逆指値注文をトレーリングストップ注文に変更します。
func (o OrderManager) DoChangeOrder(instrument string) bool {
	db := o.DB.Connect()
	distance := strconv.FormatFloat(config.GetInstance().Property.ProfitGainPrice, 'f', 5, 64)
	order, _ := o.Orders.FindLastByInstrumentAndOrder(db, instrument, enum.StopLoss)
	tradeID := strconv.Itoa(order.OrderID - 1)
	reqParam := msg.NewTradeRequest(distance, enum.Gtc)
	createRes := o.TradesApi.CreateChangeTrade(context.Background(), reqParam, tradeID)
	if createRes.ErrorCode != "" {
		log.Print("fail to change Order: ", createRes.ErrorCode, createRes.ErrorMessage)
		return false
	}
	return true
}

// 引数の注文情報を保存します。
func (o OrderManager) CreateOrder(order *domain.Orders) {
	db := o.DB.Connect()
	o.Orders.Create(db, order)
}

// 引数の取引情報を保存します。
func (o OrderManager) CreateTrade(trade *domain.Trades) {
	db := o.DB.Connect()
	o.Trades.Create(db, trade)
}

// 引数の紐付けを保存します。
func (o OrderManager) CreateBind(bind *domain.OrderTradeBinds) {
	db := o.DB.Connect()
	o.OrderTradeBinds.Create(db, bind)
}

// APIのOrderResponseをBusinessLogicのOrderに変換します。
func (o OrderManager) convertToEntityCreateOrder(res *msg.OrderResponse, timeInForce enum.TimeInForce, orderType enum.Order, distance string) *domain.Orders {
	return &domain.Orders{
		OrderID:     util.ParseInt(res.OrderFillTransaction.OrderID),
		Instrument:  res.OrderFillTransaction.Instrument,
		Units:       util.ParseFloat(res.OrderFillTransaction.Units),
		Type:        orderType,
		Price:       util.ParseFloat(res.OrderFillTransaction.FullVWAP),
		Distance:    util.ParseFloat(distance),
		Time:        res.OrderFillTransaction.Time,
		Commission:  util.ParseFloat(res.OrderFillTransaction.Commission),
		TimeInForce: timeInForce,
	}
}

// APIのOrderResponseをBusinessLogicのTradeに変換します。
func (o OrderManager) convertToEntityTrade(res *msg.OrderResponse) *domain.Trades {
	return &domain.Trades{
		TradeID:      util.ParseInt(res.OrderFillTransaction.TradeOpened.TradeID),
		Units:        util.ParseFloat(res.OrderFillTransaction.TradeOpened.Units),
		Price:        util.ParseFloat(res.OrderFillTransaction.FullVWAP),
		Instrument:   res.OrderFillTransaction.Instrument,
		InitialUnits: util.ParseFloat(res.OrderFillTransaction.Units),
		OpenTime:     res.OrderFillTransaction.Time,
	}
}

// APIのOrderResponseをBusinessLogicのOrderTradeBindsに変換します。
func (o OrderManager) convertToEntityBind(res *msg.OrderResponse) *domain.OrderTradeBinds {
	return &domain.OrderTradeBinds{
		EntryOrderID:    util.ParseInt(res.RelatedTransactionIDs[0]),
		TradeID:         util.ParseInt(res.RelatedTransactionIDs[1]),
		StopLossOrderID: util.ParseInt(res.RelatedTransactionIDs[2]),
		IsDelete:        false,
	}
}

// APIのOrderGetResponseをBusinessLogicのOrderに変換します。
func (o OrderManager) convertToEntityGetOrder(res *msg.OrderGetResponse, instrument, units string) *domain.Orders {
	return &domain.Orders{
		OrderID:     util.ParseInt(res.Order.ID),
		Instrument:  instrument,
		Units:       util.ParseFloat(units),
		Type:        enum.Order(res.Order.Type),
		Price:       util.ParseFloat(res.Order.TrailingStopValue),
		Distance:    util.ParseFloat(res.Order.Distance),
		Time:        res.Order.CreateTime,
		TimeInForce: enum.TimeInForce(res.Order.TimeInForce),
	}
}

// FIXME: 暫定対応
func (o OrderManager) convertToEntityGetOrderForStopLimit(res *msg.OrderGetResponse, instrument, units string) *domain.Orders {
	return &domain.Orders{
		OrderID:     util.ParseInt(res.Order.ID),
		Instrument:  instrument,
		Units:       util.ParseFloat(units),
		Type:        enum.Order(res.Order.Type),
		Price:       util.ParseFloat(res.Order.Price),
		Time:        res.Order.CreateTime,
		TimeInForce: enum.TimeInForce(res.Order.TimeInForce),
	}
}
