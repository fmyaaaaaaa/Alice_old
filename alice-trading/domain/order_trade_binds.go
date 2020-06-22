package domain

// 注文/取引の紐付け
type OrderTradeBinds struct {
	ID              int
	EntryOrderID    int
	TradeID         int
	StopLossOrderID int
	IsDelete        bool
}
