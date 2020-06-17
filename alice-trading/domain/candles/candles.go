package candles

import "time"

// 足データ（共通）
type Candles struct {
	Time   time.Time
	Volume float64
}

// Bidレート
type BidRate struct {
	Open  float64 `gorm:"column:open_bid"`
	Close float64 `gorm:"column:close_bid"`
	High  float64 `gorm:"column:high_bid"`
	Low   float64 `gorm:"column:low_bid"`
}

// Askレート
type AskRate struct {
	Open  float64 `gorm:"column:open_ask"`
	Close float64 `gorm:"column:close_ask"`
	High  float64 `gorm:"column:high_ask"`
	Low   float64 `gorm:"column:low_ask"`
}
