package enum

import "fmt"

type Price string

const (
	M  = Price("M")
	BA = Price("BA")
)

func (p Price) ToString() string {
	return fmt.Sprint(p)
}
