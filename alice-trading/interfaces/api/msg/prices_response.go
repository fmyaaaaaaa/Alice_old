package msg

import (
	"encoding/json"
	"time"
)

type PricesResponse struct {
	Prices []ClientPrice `json:"prices"`
}

type ClientPrice struct {
	Type        string    `json:"type"`
	Time        time.Time `json:"time"`
	Bids        []BidAsk
	Asks        []BidAsk
	CloseoutBid string `json:"closeoutBid"`
	CloseoutAsk string `json:"closeoutAsk"`
	Status      string `json:"status"`
	Tradeable   bool   `json:"tradeable"`
	Instrument  string `json:"instrument"`
}

type BidAsk struct {
	Price     string      `json:"price"`
	Liquidity json.Number `json:"liquidity"`
}
