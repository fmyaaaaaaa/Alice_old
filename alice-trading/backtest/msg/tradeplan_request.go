package msg

type TradePlanRequest struct {
	Instrument  string `json:"instrument"`
	Granularity string `json:"granularity"`
	OpenPrice   string `json:"open_price"`
	HighPrice   string `json:"high_price"`
	LowPrice    string `json:"low_price"`
	ClosePrice  string `json:"close_price"`
	Time        string `json:"time"`
}
