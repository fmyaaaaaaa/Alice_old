package oanda

import (
	"context"
	"fmt"
	"github.com/fmyaaaaaaa/Alice/alice-trading/infrastructure/config"
	"github.com/fmyaaaaaaa/Alice/alice-trading/infrastructure/logger"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/api/msg"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/api/strings"
	"log"
	"net/http"
	strings2 "strings"
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
func (p PositionsApi) GetPosition(ctx context.Context, cancel context.CancelFunc, instrument string) (*msg.PositionResponse, error) {
	defer cancel()
	strPath := fmt.Sprintf("/v3/accounts/%s/positions/%s", config.GetInstance().Api.AccountId, instrument)
	req, err := p.newRequest(ctx, "GET", strPath, nil)
	if err != nil {
		logger.LogManager().Error(err)
	}
	res, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		logger.LogManager().Error(err)
		return &msg.PositionResponse{}, err
	}
	var position msg.PositionResponse
	if err := p.decodeBody(res, &position); err != nil {
		logger.LogManager().Error(err)
	}
	return &position, nil
}

// 全銘柄のポジションを取得します。
func (p PositionsApi) GetPositions(ctx context.Context) *msg.PositionsResponse {
	strPath := fmt.Sprintf("/v3/accounts/%s/positions", config.GetInstance().Api.AccountId)
	req, err := p.newRequest(ctx, "GET", strPath, nil)
	if err != nil {
		logger.LogManager().Error(err)
	}
	res, err := p.HTTPClient.Do(req)
	if err != nil {
		logger.LogManager().Error(err)
	}
	var position msg.PositionsResponse
	if err := p.decodeBody(res, &position); err != nil {
		logger.LogManager().Error(err)
	}
	return &position
}

// 保有中のポジションを取得します。
func (p PositionsApi) GetOpenPositions(ctx context.Context) *msg.PositionsResponse {
	strPath := fmt.Sprintf("/v3/accounts/%s/openPositions", config.GetInstance().Api.AccountId)
	req, err := p.newRequest(ctx, "GET", strPath, nil)
	if err != nil {
		logger.LogManager().Error(err)
	}
	res, err := p.HTTPClient.Do(req)
	if err != nil {
		logger.LogManager().Error(err)
	}
	var position msg.PositionsResponse
	if err := p.decodeBody(res, &position); err != nil {
		logger.LogManager().Error(err)
	}
	return &position
}

func (p PositionsApi) ClosePosition(ctx context.Context, instrument string, units float64) {
	strPath := fmt.Sprintf("/v3/accounts/%s/positions/%s/close", config.GetInstance().Api.AccountId, instrument)
	var params string
	if units < 0 {
		params = `{"shortUnits": "ALL"}`
	} else {
		params = `{"longUnits": "ALL"}`
	}
	req, err := p.newRequest(ctx, "PUT", strPath, strings2.NewReader(params))
	if err != nil {
		logger.LogManager().Error(err)
	}
	res, err := p.HTTPClient.Do(req)
	if err != nil {
		logger.LogManager().Error(err)
	}
	if res.StatusCode == http.StatusOK {
		log.Println("Complete Profit Gain Order :", instrument)
	}
}
