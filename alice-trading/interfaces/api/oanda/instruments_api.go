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

// 銘柄関連のAPI
type InstrumentsApi struct {
	RootApi
}

// InstrumentsApiのコンストラクタです。
func NewInstrumentsApi() (*InstrumentsApi, error) {
	httpClient := &http.Client{}
	parsedUrl := strings.ParsedUrl(config.GetInstance().Api.Url)
	return &InstrumentsApi{RootApi{parsedUrl, httpClient}}, nil
}

// Instrumentの一覧を取得します。
func (i InstrumentsApi) GetInstruments(ctx context.Context) (*msg.InstrumentsResponse, error) {
	strPath := fmt.Sprintf("/v3/accounts/%s/instruments", config.GetInstance().Api.AccountId)
	req, err := i.newRequest(ctx, "GET", strPath, nil)
	if err != nil {
		logger.LogManager().Error(err)
		return nil, err
	}
	res, err := i.HTTPClient.Do(req)
	if err != nil {
		logger.LogManager().Error(err)
		return nil, err
	}
	var instruments msg.InstrumentsResponse
	if err := i.decodeBody(res, &instruments); err != nil {
		logger.LogManager().Error(err)
		return nil, err
	}
	return &instruments, nil
}

// 引数のinstrumentNameで指定したinstrumentを取得します。
func (i InstrumentsApi) GetInstrument(ctx context.Context, instrumentName string) (*msg.InstrumentsResponse, error) {
	strPath := fmt.Sprintf("/v3/accounts/%s/instruments", config.GetInstance().Api.AccountId)
	req, err := i.newRequest(ctx, "GET", strPath, nil)
	if err != nil {
		logger.LogManager().Error(err)
		return nil, err
	}

	params := req.URL.Query()
	params.Add("instruments", instrumentName)
	req.URL.RawQuery = params.Encode()
	res, err := i.HTTPClient.Do(req)
	if err != nil {
		logger.LogManager().Error(err)
		return nil, err
	}
	var instrument msg.InstrumentsResponse
	if err := i.decodeBody(res, &instrument); err != nil {
		logger.LogManager().Error(err)
		return nil, err
	}
	return &instrument, nil
}
