package enum

// 取引の状態
type TradeState string

const (
	Open               = TradeState("OPEN")
	Closed             = TradeState("CLOSED")
	CloseWhenTradeable = TradeState("CLOSE_WHEN_TRADEABLE")
)
