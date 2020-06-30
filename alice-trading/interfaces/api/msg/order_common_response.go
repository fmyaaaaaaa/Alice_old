package msg

import (
	"encoding/json"
	"time"
)

type Transaction struct {
	ID        string      `json:"id"`
	Time      time.Time   `json:"time"`
	UserId    json.Number `json:"userId"`
	AccountID string      `json:"accountId"`
	BatchID   string      `json:"batchId"`
	RequestID string      `json:"requestId"`
}

type OrderFillTransaction struct {
	ID                            string      `json:"id"`
	Time                          time.Time   `json:"time"`
	UserId                        int         `json:"userId"`
	AccountID                     string      `json:"accountId"`
	BatchID                       string      `json:"batchId"`
	RequestID                     string      `json:"requestId"`
	Type                          string      `json:"type"`
	OrderID                       string      `json:"orderId"`
	ClientOrderID                 string      `json:"clientOrderId"`
	Instrument                    string      `json:"instrument"`
	Units                         string      `json:"units"`
	GainQuoteHomeConversionFactor string      `json:"gainQuoteHomeConversionFactor"`
	LossQuoteHomeConversionFactor string      `json:"lossQuoteHomeConversionFactor"`
	FullVWAP                      string      `json:"fullVWAP"`
	FullPrice                     ClientPrice `json:"fullPrice"`
	Reason                        string      `json:"reason"`
	PL                            string      `json:"pl"`
	Financing                     string      `json:"financing"`
	Commission                    string      `json:"commission"`
	GuaranteedExecutionFee        string      `json:"guaranteedExecutionFee"`
	AccountBalance                string      `json:"accountBalance"`
	TradeOpened                   TradeOpen   `json:"tradeOpened"`
	HalfSpreadCost                string      `json:"halfSpreadCost"`
}

type TradeOpen struct {
	TradeID                string `json:"tradeID"`
	Units                  string `json:"units"`
	GuaranteedExecutionFee string `json:"guaranteedExecutionFee"`
	HalfSpreadCost         string `json:"halfSpreadCost"`
	InitialMarginRequired  string `json:"initialMarginRequired"`
}
