package enum

import "fmt"

// 注文タイプ
type Order string

const (
	Market           = Order("MARKET")
	Limit            = Order("LIMIT")
	Stop             = Order("STOP")
	MarketIfTouched  = Order("MARKET_IF_TOUCHED")
	TakeProfit       = Order("TAKE_PROFIT")
	StopLoss         = Order("STOP_LOSS")
	TrailingStopLoss = Order("TRAILING_STOP_LOSS")
	FixedPrice       = Order("FIXED_PRICE")
)

func (o Order) ToString() string {
	return fmt.Sprint(o)
}
