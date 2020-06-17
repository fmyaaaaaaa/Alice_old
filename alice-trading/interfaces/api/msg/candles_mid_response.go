package msg

import (
	"encoding/json"
	"time"
)

type CandlesMidResponse struct {
	Candles []CandlesMid `json:"candles"`
}

type CandlesMid struct {
	Complete bool        `json:"complete"`
	Mid      MidRate     `json:"mid"`
	Time     time.Time   `json:"time"`
	Volume   json.Number `json:"volume"`
}

type MidRate struct {
	C string `json:"c"`
	H string `json:"h"`
	L string `json:"l"`
	O string `json:"o"`
}
