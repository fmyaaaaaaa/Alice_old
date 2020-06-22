package oanda

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/fmyaaaaaaa/Alice/alice-trading/infrastructure/config"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/api/msg"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/api/strings"
	"log"
	"net/http"
)

// 注文のAPI
type OrdersApi struct {
	RootApi
}

func NewOrdersApi() *OrdersApi {
	httpClient := &http.Client{}
	parsedUrl := strings.ParsedUrl(config.GetInstance().Api.Url)
	return &OrdersApi{RootApi{HTTPClient: httpClient, URL: parsedUrl}}
}

// orderIDに一致する注文情報を取得します。
func (o OrdersApi) GetOrder(ctx context.Context, orderID string) *msg.OrderGetResponse {
	strPath := fmt.Sprintf("/v3/accounts/%s/orders/%s", config.GetInstance().Api.AccountId, orderID)
	req, err := o.newRequest(ctx, "GET", strPath, nil)
	if err != nil {
		log.Println(err)
	}
	res, err := o.HTTPClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	var order msg.OrderGetResponse
	if err := o.decodeBody(res, &order); err != nil {
		log.Println(err)
	}
	return &order
}

// OrderRequestに応じて新規注文を実行します。
func (o OrdersApi) CreateNewOrder(ctx context.Context, reqParam *msg.OrderRequest) (*msg.OrderResponse, *msg.OrderErrorResponse) {
	strPath := fmt.Sprintf("/v3/accounts/%s/orders", config.GetInstance().Api.AccountId)
	req, err := o.newRequest(ctx, "POST", strPath, createBody(reqParam))
	if err != nil {
		log.Print(err)
	}

	res, err := o.HTTPClient.Do(req)
	if err != nil {
		log.Print(err)
	}

	if res.StatusCode == http.StatusCreated {
		var orderResponse msg.OrderResponse
		if err := o.decodeBody(res, &orderResponse); err != nil {
			log.Println(err)
		}
		return &orderResponse, nil
	} else {
		var orderErrorResponse msg.OrderErrorResponse
		if err := o.decodeBody(res, &orderErrorResponse); err != nil {
			log.Println(err)
		}
		return nil, &orderErrorResponse
	}
}

// OrderRequestに応じて注文内容を変更します。
// 引数に注文IDを指定します。
func (o OrdersApi) CreateChangeOrder(ctx context.Context, reqParam *msg.OrderRequest, orderID string) (*msg.OrderResponse, *msg.OrderErrorResponse) {
	strPath := fmt.Sprintf("/v3/accounts/%s/orders/%s", config.GetInstance().Api.AccountId, orderID)
	req, err := o.newRequest(ctx, "PUT", strPath, createBody(reqParam))
	if err != nil {
		log.Print(err)
	}

	res, err := o.HTTPClient.Do(req)
	if err != nil {
		log.Print(err)
	}

	if res.StatusCode == http.StatusCreated {
		var orderResponse msg.OrderResponse
		if err := o.decodeBody(res, &orderResponse); err != nil {
			log.Println(err)
		}
		return &orderResponse, nil
	} else {
		var orderErrorResponse msg.OrderErrorResponse
		if err := o.decodeBody(res, &orderErrorResponse); err != nil {
			log.Println(err)
		}
		return nil, &orderErrorResponse
	}
}

func createBody(reqParam *msg.OrderRequest) *bytes.Buffer {
	jsonByte, err := json.Marshal(reqParam)
	if err != nil {
		log.Println("fail to parse json ", err)
	}
	return bytes.NewBuffer(jsonByte)
}
