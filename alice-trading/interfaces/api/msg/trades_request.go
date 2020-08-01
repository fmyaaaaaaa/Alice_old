package msg

import "github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"

type TradesRequest struct {
	TrailingStopLossDetails TrailingStopLossDetails `json:"trailingStopLoss,omitempty"`
}

// TrailingStopのRequestを生成します。
func NewTradeRequest(distance string, timeInForce enum.TimeInForce) *TradesRequest {
	return &TradesRequest{
		TrailingStopLossDetails: TrailingStopLossDetails{
			Distance:    distance,
			TimeInForce: timeInForce,
		}}
}
