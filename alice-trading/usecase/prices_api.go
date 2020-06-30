package usecase

import (
	"context"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/api/msg"
)

// レートのAPI
type PricesApi interface {
	GetPrices(ctx context.Context, instrument string) *msg.PricesResponse
}
