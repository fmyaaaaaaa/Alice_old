package oanda

import (
	"context"
	"encoding/json"
	"github.com/fmyaaaaaaa/Alice/alice-trading/infrastructure/config"
	"io"
	"net/http"
	"net/url"
	"path"
)

type RootApi struct {
	URL        *url.URL
	HTTPClient *http.Client
}

func (r RootApi) newRequest(ctx context.Context, method, strPath string, body io.Reader) (*http.Request, error) {
	u := *r.URL
	u.Path = path.Join(r.URL.Path, strPath)
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	for key, value := range r.createHeader() {
		req.Header.Set(key, value)
	}

	return req, nil
}

func (r RootApi) decodeBody(res *http.Response, out interface{}) error {
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	return decoder.Decode(out)
}

func (r RootApi) decodeBodyForStreaming(res *http.Response, out interface{}) error {
	decoder := json.NewDecoder(res.Body)
	return decoder.Decode(out)
}

func (r RootApi) createHeader() map[string]string {
	return map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + config.GetInstance().Api.AccessToken,
	}
}
