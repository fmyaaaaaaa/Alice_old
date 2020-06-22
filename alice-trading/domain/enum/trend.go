package enum

import "fmt"

// トレンド
type Trend string

const (
	UpTrend   = Trend("UPTREND")
	DownTrend = Trend("DOWNTREND")
	Range     = Trend("RANGE")
)

func (t Trend) ToString() string {
	return fmt.Sprint(t)
}
