package usecase

import (
	"context"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/api/msg"
)

// 注文のAPI
type OrdersApi interface {
	GetOrder(ctx context.Context, orderID string) *msg.OrderGetResponse
	CreateNewOrder(ctx context.Context, reqParam *msg.OrderRequest) (*msg.OrderResponse, *msg.OrderErrorResponse)
	CreateChangeOrder(ctx context.Context, reqParam *msg.OrderRequest, orderID string) (*msg.OrderResponse, *msg.OrderErrorResponse)
}
