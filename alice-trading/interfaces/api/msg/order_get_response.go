package msg

import (
	"time"
)

type OrderGetResponse struct {
	Order             CommonOrderResponse `json:"order"`
	LastTransactionID string              `json:"lastTransactionID"`
}

type CommonOrderResponse struct {
	ID                string    `json:"id"`
	CreateTime        time.Time `json:"createTime"`
	Type              string    `json:"type"`
	TradeID           string    `json:"tradeId"`
	Distance          string    `json:"distance"`
	TimeInForce       string    `json:"timeInForce"`
	TriggerCondition  string    `json:"triggerCondition"`
	State             string    `json:"state"`
	TrailingStopValue string    `json:"trailingStopValue"`
}
