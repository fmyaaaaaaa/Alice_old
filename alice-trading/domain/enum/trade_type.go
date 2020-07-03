package enum

// トレーディング種別
type TradeType string

const (
	DayTrade   = TradeType("DAY_TRADE")
	SwingTrade = TradeType("SWING_TRADE")
)
