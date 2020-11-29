package oanda

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/fmyaaaaaaa/Alice/alice-trading/infrastructure/config"
	"github.com/fmyaaaaaaa/Alice/alice-trading/infrastructure/logger"
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

// 取引情報を取得します。
func (t TradesApi) GetTrades(ctx context.Context) *msg.TradesResponse {
	strPath := fmt.Sprintf("/v3/accounts/%s/trades", config.GetInstance().Api.AccountId)
	req, err := t.newRequest(ctx, "GET", strPath, nil)
	if err != nil {
		logger.LogManager().Error(err)
	}
	res, err := t.HTTPClient.Do(req)
	if err != nil {
		logger.LogManager().Error(err)
	}
	var trade msg.TradesResponse
	if err := t.decodeBody(res, &trade); err != nil {
		logger.LogManager().Error(err)
	}
	return &trade
}

// 取引中の取引情報を変更します。
func (t TradesApi) CreateChangeTrade(ctx context.Context, reqParam *msg.TradesRequest, tradeID string) *msg.TradesChangeResponse {
	strPath := fmt.Sprintf("/v3/accounts/%s/trades/orders", config.GetInstance().Api.AccountId)
	req, err := t.newRequest(ctx, "PUT", strPath, createBodyTradesRequest(reqParam))
	if err != nil {
		log.Print(err)
	}

	res, err := t.HTTPClient.Do(req)
	if err != nil {
		log.Print(err)
	}

	var tradeChangeResponse msg.TradesChangeResponse
	if err := t.decodeBody(res, &tradeChangeResponse); err != nil {
		logger.LogManager().Error(err)
	}
	return &tradeChangeResponse
}

func createBodyTradesRequest(reqParam *msg.TradesRequest) *bytes.Buffer {
	jsonByte, err := json.Marshal(reqParam)
	if err != nil {
		log.Println("fail to parse json ", err)
	}
	return bytes.NewBuffer(jsonByte)
}
