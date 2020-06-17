package msg

import (
	"encoding/json"
	"time"
)

type CandlesBidAskResponse struct {
	Candles []CandlesBidAsk `json:"candles"`
}

type CandlesBidAsk struct {
	Ask      BidAskRate  `json:"ask"`
	Bid      BidAskRate  `json:"bid"`
	Complete bool        `json:"complete"`
	Time     time.Time   `json:"time"`
	Volume   json.Number `json:"volume"`
}

type BidAskRate struct {
	C string `json:"c"`
	H string `json:"h"`
	L string `json:"l"`
	O string `json:"o"`
}
