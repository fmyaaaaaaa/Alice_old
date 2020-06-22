package enum

import "fmt"

// 足データの陽線/陰線
type Line string

const (
	Positive = Line("POSITIVE")
	Negative = Line("NEGATIVE")
)

func (l Line) ToString() string {
	return fmt.Sprint(l)
}
