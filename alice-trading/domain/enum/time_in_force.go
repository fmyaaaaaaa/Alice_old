package enum

import "fmt"

type TimeInForce string

const (
	Gtc = TimeInForce("GTC")
	Gtd = TimeInForce("GTD")
	Gfd = TimeInForce("GFD")
	Fok = TimeInForce("FOK")
	Ioc = TimeInForce("IOC")
)

func (t TimeInForce) ToString() string {
	return fmt.Sprint(t)
}
