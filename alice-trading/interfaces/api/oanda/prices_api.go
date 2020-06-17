package oanda

import (
	"context"
	"fmt"
	"github.com/fmyaaaaaaa/Alice/alice-trading/infrastructure/config"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/api/msg"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/api/strings"
	"net/http"
)

// リアルタイムレートのAPI
type PricesApi struct {
	RootApi
}

func NewPricesApi() (*PricesApi, error) {
	httpClient := &http.Client{}
	parsedUrl := strings.ParsedUrl("https://stream-fxpractice.oanda.com")
	return &PricesApi{RootApi{URL: parsedUrl, HTTPClient: httpClient}}, nil
}

// リアルタイムレートを取得します。（HttpStreaming）
func (p PricesApi) StreamingPrices(ctx context.Context, instrumentName string) {
	strPath := fmt.Sprintf("/v3/accounts/%s/pricing/stream", config.GetInstance().Api.AccountId)
	req, err := p.newRequest(ctx, "GET", strPath, nil)
	if err != nil {
		panic(err)
	}

	params := req.URL.Query()
	params.Add("instruments", instrumentName)
	req.URL.RawQuery = params.Encode()

	res, err := p.HTTPClient.Do(req)
	if err != nil {
		panic(err)
	}

	for {
		var prices msg.PricesResponse
		if err := p.decodeBodyForStreaming(res, &prices); err != nil {
			panic(err)
		}
	}
}
