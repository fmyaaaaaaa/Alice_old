package usecase

import (
	"context"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/api/msg"
)

// アカウントのAPI
type AccountsApi interface {
	GetAccountSummary(ctx context.Context, cancel context.CancelFunc) (*msg.AccountSummaryResponse, error)
}
