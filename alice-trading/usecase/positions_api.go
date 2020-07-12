package usecase

import (
	"context"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/api/msg"
)

// ポジションのAPI
type PositionsApi interface {
	GetPosition(ctx context.Context, cancel context.CancelFunc, instrument string) (*msg.PositionResponse, error)
	GetPositions(ctx context.Context) *msg.PositionsResponse
	GetOpenPositions(ctx context.Context) *msg.PositionsResponse
	ClosePosition(ctx context.Context, instrument string, units float64)
}
