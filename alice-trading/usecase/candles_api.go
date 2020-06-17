package usecase

import (
	"context"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/api/msg"
)

// 足データのAPI
type CandlesApi interface {
	GetCandleMid(ctx context.Context, instrumentName string, count int, granularity enum.Granularity) (*msg.CandlesMidResponse, error)
	GetCandleBidAsk(ctx context.Context, instrumentName string, count int, granularity enum.Granularity) (*msg.CandlesBidAskResponse, error)
}
