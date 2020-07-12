package msg

import (
	"log"
	"strconv"
)

type BuySell string

const (
	Buy  = BuySell("BUY")
	Sell = BuySell("SELL")
)

type TradePlanResponse struct {
	Status   int     `json:"status"`
	Result   string  `json:"result"`
	Units    string  `json:"units"`
	Distance string  `json:"distance"`
	BuySell  BuySell `json:"buy_sell"`
	IsOrder  bool    `json:"is_order"`
}

func NewTradePlanResponse(status int, result, units string, distance float64, isOrder bool) TradePlanResponse {
	var buySell BuySell
	switch {
	case parseFloat(units) > 0:
		buySell = Buy
	case parseFloat(units) < 0:
		buySell = Sell
	}
	strDistance := strconv.FormatFloat(distance, 'f', 5, 64)
	return TradePlanResponse{
		Status:   status,
		Result:   result,
		Units:    units,
		Distance: strDistance,
		BuySell:  buySell,
		IsOrder:  isOrder,
	}
}

// 文字列をfloat64に変換します。
func parseFloat(str string) float64 {
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Print("fail to parse float64")
	}
	return f
}
