package oanda

import (
	"context"
	"fmt"
	"github.com/fmyaaaaaaa/Alice/alice-trading/infrastructure/config"
	"github.com/fmyaaaaaaa/Alice/alice-trading/infrastructure/logger"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/api/msg"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/api/strings"
	"net/http"
)

// レートのAPI
type PricesApi struct {
	RootApi
}

func NewPricesApi() *PricesApi {
	httpClient := &http.Client{}
	parsedUrl := strings.ParsedUrl(config.GetInstance().Api.Url)
	return &PricesApi{RootApi{URL: parsedUrl, HTTPClient: httpClient}}
}

// 指定した銘柄のレートを取得します。
func (p PricesApi) GetPrices(ctx context.Context, instrument string) *msg.PricesResponse {
	strPath := fmt.Sprintf("/v3/accounts/%s/pricing", config.GetInstance().Api.AccountId)
	req, err := p.newRequest(ctx, "GET", strPath, nil)
	params := req.URL.Query()
	params.Add("instruments", instrument)
	req.URL.RawQuery = params.Encode()
	if err != nil {
		logger.LogManager().Error(err)
	}
	res, err := p.HTTPClient.Do(req)
	if err != nil {
		logger.LogManager().Error(err)
	}
	var price msg.PricesResponse
	if err := p.decodeBody(res, &price); err != nil {
		logger.LogManager().Error(err)
	}
	return &price
}
