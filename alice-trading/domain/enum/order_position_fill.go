package enum

import "fmt"

type OrderPositionFill string

const (
	OpenOnly                 = OrderPositionFill("OPEN_ONLY")
	ReduceFirst              = OrderPositionFill("REDUCE_FIRST")
	ReduceOnly               = OrderPositionFill("REDUCE_ONLY")
	DefaultOrderPositionFIll = OrderPositionFill("DEFAULT")
)

func (o OrderPositionFill) ToString() string {
	return fmt.Sprint(o)
}
