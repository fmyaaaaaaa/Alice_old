package enum

type PriceType int

const (
	Mid PriceType = iota
	BidAsk
)

func (p PriceType) ConvertToParam() string {
	switch p {
	case Mid:
		return "M"
	case BidAsk:
		return "BA"
	default:
		return ""
	}
}
