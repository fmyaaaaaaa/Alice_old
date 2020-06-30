package msg

import (
	"time"
)

type TradesResponse struct {
	Trades            []Trade `json:"trades"`
	LastTransactionID string  `json:"lastTransactionID"`
}

type Trade struct {
	ID                    string                  `json:"id"`
	Instrument            string                  `json:"instrument"`
	Price                 string                  `json:"price"`
	OpenTime              time.Time               `json:"openTime"`
	State                 string                  `json:"state"`
	InitialUnits          string                  `json:"initialUnits"`
	InitialMarginRequired string                  `json:"initialMarginRequired"`
	CurrentUnits          string                  `json:"currentUnits"`
	RealizedPl            string                  `json:"realizedPL"`
	UnrealizedPl          string                  `json:"unrealizedPL"`
	MarginUsed            string                  `json:"marginUsed"`
	AverageClosePrice     string                  `json:"averageClosePrice"`
	ClosingTransactionIDs []string                `json:"closingTransactionIDs"`
	Financing             string                  `json:"financing"`
	DividendAdjustment    string                  `json:"dividendAdjustment"`
	CloseTime             time.Time               `json:"closeTime"`
	TrailingStopLossOrder []TrailingStopLossOrder `json:"trailingStopLossOrder"`
}

type TrailingStopLossOrder struct {
	ID                   string    `json:"id"`
	CreateTime           time.Time `json:"createTime"`
	State                string    `json:"state"`
	Type                 string    `json:"type"`
	TradeID              string    `json:"tradeID"`
	ClientTradeID        string    `json:"clientTradeID"`
	Distance             string    `json:"distance"`
	TimeInForce          string    `json:"timeInForce"`
	GtdTime              time.Time `json:"gtdTime"`
	TriggerCondition     string    `json:"triggerCondition"`
	TrailingStopValue    string    `json:"trailingStopValue"`
	FillingTransactionID string    `json:"fillingTransactionID"`
	FilledTime           time.Time `json:"filledTime"`
	TradeOpenedID        string    `json:"tradeOpenedID"`
}
