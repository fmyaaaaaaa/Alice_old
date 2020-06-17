package enum

import "fmt"

// 足データの陽線/陰線
type Line string

const (
	POSITIVE = Line("POSITIVE")
	NEGATIVE = Line("NEGATIVE")
)

func (l Line) ToString() string {
	return fmt.Sprint(l)
}
