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

// アカウントのAPI
type AccountsApi struct {
	RootApi
}

func NewAccountApi() *AccountsApi {
	httpClient := &http.Client{}
	parsedUrl := strings.ParsedUrl(config.GetInstance().Api.Url)
	return &AccountsApi{RootApi{URL: parsedUrl, HTTPClient: httpClient}}
}

// アカウントのサマリーを取得します。
func (a AccountsApi) GetAccountSummary(ctx context.Context) *msg.AccountSummaryResponse {
	strPath := fmt.Sprintf("/v3/accounts/%s/summary", config.GetInstance().Api.AccountId)
	req, err := a.newRequest(ctx, "GET", strPath, nil)
	if err != nil {
		log.Println(err)
	}
	res, err := a.HTTPClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	var accountSummary msg.AccountSummaryResponse
	if err := a.decodeBody(res, &accountSummary); err != nil {
		log.Println(err)
	}
	return &accountSummary
}
