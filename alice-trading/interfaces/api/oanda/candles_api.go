package oanda

import (
	"context"
	"fmt"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"github.com/fmyaaaaaaa/Alice/alice-trading/infrastructure/config"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/api/msg"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/api/strings"
	"net/http"
	"strconv"
)

// 足データのAPI
type CandlesApi struct {
	RootApi
}

func NewCandlesApi() *CandlesApi {
	httpClient := &http.Client{}
	parsedUrl := strings.ParsedUrl(config.GetInstance().Api.Url)
	return &CandlesApi{RootApi{HTTPClient: httpClient, URL: parsedUrl}}
}

// Midの足データを取得します。
// 銘柄名、本数、足種を指定します。
func (c CandlesApi) GetCandleMid(ctx context.Context, instrumentName string, count int, granularity enum.Granularity) (*msg.CandlesMidResponse, error) {
	strPath := fmt.Sprintf("/v3/instruments/%s/candles", instrumentName)
	req, err := c.newRequest(ctx, "GET", strPath, nil)
	if err != nil {
		panic(err)
	}
	createCandleParam(req, strconv.Itoa(count), enum.Mid.ConvertToParam(), granularity.ToString())
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		panic(err)
	}

	var candles msg.CandlesMidResponse
	if err := c.decodeBody(res, &candles); err != nil {
		panic(err)
	}
	return &candles, nil
}

// BidAskの足データを取得します。
// 銘柄名、本数、足種を指定します。
func (c CandlesApi) GetCandleBidAsk(ctx context.Context, instrumentName string, count int, granularity enum.Granularity) (*msg.CandlesBidAskResponse, error) {
	strPath := fmt.Sprintf("/v3/instruments/%s/candles", instrumentName)
	req, err := c.newRequest(ctx, "GET", strPath, nil)
	if err != nil {
		panic(err)
	}
	createCandleParam(req, strconv.Itoa(count), enum.BidAsk.ConvertToParam(), granularity.ToString())
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		panic(err)
	}

	var candles msg.CandlesBidAskResponse
	if err := c.decodeBody(res, &candles); err != nil {
		panic(err)
	}
	return &candles, nil
}

// 足データ取得のパラメータを設定します。
func createCandleParam(req *http.Request, count string, price string, granularity string) {
	params := req.URL.Query()
	params.Add("count", count)
	params.Add("price", price)
	params.Add("granularity", granularity)
	req.URL.RawQuery = params.Encode()
}
