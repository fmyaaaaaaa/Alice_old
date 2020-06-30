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

// ポジションのAPI
type PositionsApi struct {
	RootApi
}

func NewPositionsApi() *PositionsApi {
	httpClient := &http.Client{}
	parsedUrl := strings.ParsedUrl(config.GetInstance().Api.Url)
	return &PositionsApi{RootApi{URL: parsedUrl, HTTPClient: httpClient}}
}

// 銘柄を指定してポジションを取得します。
func (p PositionsApi) GetPosition(ctx context.Context, instrument string) *msg.PositionResponse {
	strPath := fmt.Sprintf("/v3/accounts/%s/positions/%s", config.GetInstance().Api.AccountId, instrument)
	req, err := p.newRequest(ctx, "GET", strPath, nil)
	if err != nil {
		log.Println(err)
	}
	res, err := p.HTTPClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	var position msg.PositionResponse
	if err := p.decodeBody(res, &position); err != nil {
		log.Println(err)
	}
	return &position
}

// 全銘柄のポジションを取得します。
func (p PositionsApi) GetPositions(ctx context.Context) *msg.PositionsResponse {
	strPath := fmt.Sprintf("/v3/accounts/%s/positions", config.GetInstance().Api.AccountId)
	req, err := p.newRequest(ctx, "GET", strPath, nil)
	if err != nil {
		log.Println(err)
	}
	res, err := p.HTTPClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	var position msg.PositionsResponse
	if err := p.decodeBody(res, &position); err != nil {
		log.Println(err)
	}
	return &position
}

// 保有中のポジションを取得します。
func (p PositionsApi) GetOpenPositions(ctx context.Context) *msg.PositionsResponse {
	strPath := fmt.Sprintf("/v3/accounts/%s/openPositions", config.GetInstance().Api.AccountId)
	req, err := p.newRequest(ctx, "GET", strPath, nil)
	if err != nil {
		log.Println(err)
	}
	res, err := p.HTTPClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	var position msg.PositionsResponse
	if err := p.decodeBody(res, &position); err != nil {
		log.Println(err)
	}
	return &position
}
