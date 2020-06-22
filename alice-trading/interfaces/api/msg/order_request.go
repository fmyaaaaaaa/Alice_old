package msg

import "github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"

type OrderRequest struct {
	Order OrderRequestParam `json:"order"`
}

type OrderRequestParam struct {
	Type                   enum.Order                 `json:"type,omitempty"`
	Instrument             string                     `json:"instrument,omitempty"`
	Units                  string                     `json:"units,omitempty"`
	Price                  string                     `json:"price,omitempty"`
	TimeInForce            enum.TimeInForce           `json:"timeInForce,omitempty"`
	PriceBound             string                     `json:"priceBound,omitempty"`
	PositionFill           enum.OrderPositionFill     `json:"positionFill,omitempty"`
	TriggerCondition       enum.OrderTriggerCondition `json:"triggerCondition,omitempty"`
	TakeProfitOnFill       *TakeProfitDetails         `json:"takeProfitOnFill,omitempty"`
	StopLossOnFill         *StopLossDetails           `json:"stopLossOnFill,omitempty"`
	TrailingStopLossOnFill *TrailingStopLossDetails   `json:"trailingStopLossOnFill,omitempty"`
	// -- ProfitOrder --
	TradeID       string `json:"tradeID,omitempty"`
	ClientTradeID string `json:"clientTradeID,omitempty"`
	Distance      string `json:"distance,omitempty"`
}

// MarketOrderのRequestを生成します。
func NewMarketOrderRequest(instrument string, units string, priceBound string, orderPositionFill enum.OrderPositionFill, takeProfitDetails *TakeProfitDetails, stopLossDetails *StopLossDetails, trailingStopLossDetails *TrailingStopLossDetails) *OrderRequest {
	return &OrderRequest{
		Order: OrderRequestParam{
			Type:                   enum.Market,
			Instrument:             instrument,
			Units:                  units,
			TimeInForce:            enum.Fok,
			PriceBound:             priceBound,
			PositionFill:           orderPositionFill,
			TakeProfitOnFill:       takeProfitDetails,
			StopLossOnFill:         stopLossDetails,
			TrailingStopLossOnFill: trailingStopLossDetails,
		}}
}

// LimitOrderのRequestを生成します。
func NewLimitOrderRequest(instrument string, units string, price string, timeInForce enum.TimeInForce, orderPositionFill enum.OrderPositionFill, orderTriggerCondition enum.OrderTriggerCondition, takeProfitDetails *TakeProfitDetails, stopLossDetails *StopLossDetails, trailingStopLossDetails *TrailingStopLossDetails) *OrderRequest {
	return &OrderRequest{
		Order: OrderRequestParam{
			Type:                   enum.Limit,
			Instrument:             instrument,
			Units:                  units,
			Price:                  price,
			TimeInForce:            timeInForce,
			PositionFill:           orderPositionFill,
			TriggerCondition:       orderTriggerCondition,
			TakeProfitOnFill:       takeProfitDetails,
			StopLossOnFill:         stopLossDetails,
			TrailingStopLossOnFill: trailingStopLossDetails,
		}}
}

// StopOrderのRequestを生成します。
func NewStopOrderRequest(instrument string, units string, price string, priceBound string, timeInForce enum.TimeInForce, orderPositionFill enum.OrderPositionFill, orderTriggerCondition enum.OrderTriggerCondition, takeProfitDetails *TakeProfitDetails, stopLossDetails *StopLossDetails, trailingStopLossDetails *TrailingStopLossDetails) *OrderRequest {
	return &OrderRequest{
		Order: OrderRequestParam{
			Type:                   enum.Stop,
			Instrument:             instrument,
			Units:                  units,
			Price:                  price,
			TimeInForce:            timeInForce,
			PriceBound:             priceBound,
			PositionFill:           orderPositionFill,
			TriggerCondition:       orderTriggerCondition,
			TakeProfitOnFill:       takeProfitDetails,
			StopLossOnFill:         stopLossDetails,
			TrailingStopLossOnFill: trailingStopLossDetails,
		}}
}

// TakeProfitOrderのRequestを生成します。
func NewTakeProfitOrderRequest(tradeID string, clientTradeID string, price string, timeInForce enum.TimeInForce, orderTriggerCondition enum.OrderTriggerCondition) *OrderRequest {
	return &OrderRequest{
		Order: OrderRequestParam{
			Type:             enum.TakeProfit,
			Price:            price,
			TimeInForce:      timeInForce,
			TradeID:          tradeID,
			ClientTradeID:    clientTradeID,
			TriggerCondition: orderTriggerCondition,
		}}
}

// StopLossOrderのRequestを生成します。
func NewStopLossOrderRequest(tradeID string, clientTradeID string, price string, distance string, timeInForce enum.TimeInForce, orderTriggerCondition enum.OrderTriggerCondition) *OrderRequest {
	return &OrderRequest{
		Order: OrderRequestParam{
			Type:             enum.StopLoss,
			Price:            price,
			TimeInForce:      timeInForce,
			TriggerCondition: orderTriggerCondition,
			TradeID:          tradeID,
			ClientTradeID:    clientTradeID,
			Distance:         distance,
		}}
}

// TrailingStopLossOrderのRequestを生成します。
func NewTrailingStopLossOrderRequest(tradeID string, clientTradeID string, distance string, timeInForce enum.TimeInForce, orderTriggerCondition enum.OrderTriggerCondition) *OrderRequest {
	return &OrderRequest{
		Order: OrderRequestParam{
			Type:             enum.TrailingStopLoss,
			TimeInForce:      timeInForce,
			TriggerCondition: orderTriggerCondition,
			TradeID:          tradeID,
			ClientTradeID:    clientTradeID,
			Distance:         distance,
		}}
}

// ChangeRequest
func NewChangeOrderRequest(orderType enum.Order, tradeID string, price string, distance string, timeInForce enum.TimeInForce) *OrderRequest {
	return &OrderRequest{Order: OrderRequestParam{
		Type:        orderType,
		Price:       price,
		TimeInForce: timeInForce,
		TradeID:     tradeID,
		Distance:    distance,
	}}
}
