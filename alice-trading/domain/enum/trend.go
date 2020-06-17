package enum

import "fmt"

// トレンド
type Trend string

const (
	UPTREND   = Trend("UPTREND")
	DOWNTREND = Trend("DOWNTREND")
	RANGE     = Trend("RANGE")
)

func (t Trend) ToString() string {
	return fmt.Sprint(t)
}
