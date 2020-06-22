package domain

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"time"
)

// 注文
type Orders struct {
	ID          int
	OrderID     int `toml:"column:order_id"`
	Instrument  string
	Units       float64
	Type        enum.Order
	Price       float64
	Distance    float64
	Time        time.Time
	Commission  float64
	TimeInForce enum.TimeInForce
}
