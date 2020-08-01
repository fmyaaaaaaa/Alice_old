package usecase

import (
	"context"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/api/msg"
)

// 取引のAPI
type TradesApi interface {
	GetTrades(ctx context.Context) *msg.TradesResponse
	CreateChangeTrade(ctx context.Context, reqParam *msg.TradesRequest, tradeID string) *msg.TradesChangeResponse
}
