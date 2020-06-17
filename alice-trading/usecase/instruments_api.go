package usecase

import (
	"context"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/api/msg"
)

// 銘柄関連のAPI
type InstrumentsApi interface {
	GetInstruments(ctx context.Context) (*msg.InstrumentsResponse, error)
	GetInstrument(ctx context.Context, instrumentName string) (*msg.InstrumentsResponse, error)
}
