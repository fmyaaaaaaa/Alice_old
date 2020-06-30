package oanda

import (
	"context"
	"fmt"
	"github.com/fmyaaaaaaa/Alice/alice-trading/infrastructure/config"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/api/msg"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/api/strings"
	"log"
	"net/http"
)

// 取引のAPI
type TradesApi struct {
	RootApi
}

func NewTradeApi() *TradesApi {
	httpClient := &http.Client{}
	parsedUrl := strings.ParsedUrl(config.GetInstance().Api.Url)
	return &TradesApi{RootApi{URL: parsedUrl, HTTPClient: httpClient}}
}

func (t TradesApi) GetTrades(ctx context.Context) *msg.TradesResponse {
	strPath := fmt.Sprintf("/v3/accounts/%s/trades", config.GetInstance().Api.AccountId)
	req, err := t.newRequest(ctx, "GET", strPath, nil)
	if err != nil {
		log.Println(err)
	}
	res, err := t.HTTPClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	var trade msg.TradesResponse
	if err := t.decodeBody(res, &trade); err != nil {
		log.Println(err)
	}
	return &trade
}
