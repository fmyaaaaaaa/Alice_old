package usecase

import (
	"context"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/api/msg"
)

// ポジションのAPI
type PositionsApi interface {
	GetPosition(ctx context.Context, instrument string) *msg.PositionResponse
	GetPositions(ctx context.Context) *msg.PositionsResponse
	GetOpenPositions(ctx context.Context) *msg.PositionsResponse
}
