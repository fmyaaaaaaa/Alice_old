package usecase

import (
	"context"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/api/msg"
)

// 取引のAPI
type TradesApi interface {
	GetTrades(ctx context.Context) *msg.TradesResponse
}
