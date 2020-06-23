package domain

// スイングの高値/安値
type SwingHighLowPrice struct {
	ID        int
	SwingID   int
	HighPrice float64
	LowPrice  float64
}

func NewSwingHighLowPrice(swingID int, highPrice, lowPrice float64) *SwingHighLowPrice {
	return &SwingHighLowPrice{
		SwingID:   swingID,
		HighPrice: highPrice,
		LowPrice:  lowPrice,
	}
}
