package msg

import "github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"

type TakeProfitDetails struct {
	Price       string           `json:"price,omitempty"`
	TimeInForce enum.TimeInForce `json:"timeInForce,omitempty"`
}

type StopLossDetails struct {
	Price       string           `json:"price,omitempty"`
	Distance    string           `json:"distance,omitempty"`
	TimeInForce enum.TimeInForce `json:"timeInForce,omitempty"`
}

type TrailingStopLossDetails struct {
	Distance    string           `json:"distance,omitempty"`
	TimeInForce enum.TimeInForce `json:"timeInForce,omitempty"`
}
