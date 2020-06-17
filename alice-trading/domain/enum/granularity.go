package enum

import "fmt"

// 足種
type Granularity string

const (
	S5  = Granularity("S5")
	S10 = Granularity("S10")
	S15 = Granularity("S15")
	S30 = Granularity("S30")
	M1  = Granularity("M1")
	M2  = Granularity("M2")
	M4  = Granularity("M4")
	M5  = Granularity("M5")
	M10 = Granularity("M10")
	M15 = Granularity("M15")
	M30 = Granularity("M30")
	H1  = Granularity("H1")
	H2  = Granularity("H2")
	H3  = Granularity("H3")
	H4  = Granularity("H4")
	H6  = Granularity("H6")
	H8  = Granularity("H8")
	H12 = Granularity("H12")
	D   = Granularity("D")
	W   = Granularity("W")
	M   = Granularity("M")
)

func (g Granularity) ToString() string {
	return fmt.Sprint(g)
}
