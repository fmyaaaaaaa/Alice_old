package money

import (
	"context"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/api/msg"
	"github.com/fmyaaaaaaa/Alice/alice-trading/usecase"
	"github.com/fmyaaaaaaa/Alice/alice-trading/usecase/util"
)

// アカウント情報管理
type AccountManager struct {
	DB           usecase.DBRepository
	Accounts     usecase.AccountRepository
	Trades       usecase.TradesRepository
	Positions    usecase.PositionsRepository
	AccountsApi  usecase.AccountsApi
	TradesApi    usecase.TradesApi
	PositionsApi usecase.PositionsApi
}

// アカウント情報を更新し、更新後のアカウント情報を返却します。
func (a AccountManager) UpdateAccountInformation() domain.Accounts {
	// アカウント情報の最新化
	res := a.AccountsApi.GetAccountSummary(context.Background())
	return a.updateAccount(res)
}

// アカウント情報を取得します。
// ユーザー管理は考慮不要のため、IDを1で指定しています。
func (a AccountManager) GetAccount() domain.Accounts {
	DB := a.DB.Connect()
	return a.Accounts.FindByID(DB, 1)
}

// ポジションの保有状況を取得します。
func (a AccountManager) HasPosition(instrument string) (bool, domain.Positions) {
	position := a.GetPosition(instrument)
	if position.Units != 0 {
		return true, position
	}
	return false, position
}

// ポジション情報を新規作成または更新します。
func (a AccountManager) CreateOrUpdatePosition(instrument string) domain.Positions {
	DB := a.DB.Connect()
	res := a.PositionsApi.GetPosition(context.Background(), instrument)
	position := a.convertToEntityPosition(res)
	a.Positions.CreateOrUpdate(DB, &position)
	return position
}

// 指定した銘柄のポジションを更新し、返却します。
func (a AccountManager) UpdatePositionInformation(instrument string) domain.Positions {
	res := a.PositionsApi.GetPosition(context.Background(), instrument)
	position := a.convertToEntityPosition(res)
	a.updatePosition(&position)
	return position
}

// 全銘柄のポジションを更新し、返却します。
func (a AccountManager) UpdatePositionsInformation() *[]domain.Positions {
	res := a.PositionsApi.GetPositions(context.Background())
	positions := a.convertToEntityPositions(res)
	a.updateBulkPositions(positions)
	return positions
}

// 指定した銘柄のポジションをクローズします。
func (a AccountManager) ClosePosition(instrument string, units float64) {
	a.PositionsApi.ClosePosition(context.Background(), instrument, units)
}

// アカウント情報を更新します。
func (a AccountManager) updateAccount(res *msg.AccountSummaryResponse) domain.Accounts {
	DB := a.DB.Connect()
	return a.Accounts.Update(DB, a.createUpdateParamsAccount(res))
}

// 指定した銘柄のポジションを取得します。
func (a AccountManager) GetPosition(instrument string) domain.Positions {
	DB := a.DB.Connect()
	return a.Positions.FindByInstrument(DB, instrument)
}

// ポジション情報を更新します。
func (a AccountManager) updatePosition(position *domain.Positions) {
	DB := a.DB.Connect()
	a.Positions.Update(DB, position)
}

// ポジション情報を一括で更新します。
func (a AccountManager) updateBulkPositions(positions *[]domain.Positions) {
	DB := a.DB.Connect()
	a.Positions.BulkUpdate(DB, positions)
}

// APIのレスポンスからアカウント情報の更新用パラメータを作成します。
func (a AccountManager) createUpdateParamsAccount(res *msg.AccountSummaryResponse) map[string]interface{} {
	params := map[string]interface{}{
		"margin_rate":                    util.ParseFloat(res.Account.MarginRate),
		"balance":                        util.ParseFloat(res.Account.Balance),
		"open_trade_count":               util.ParseFloat(res.Account.OpenTradeCount.String()),
		"open_position_count":            util.ParseFloat(res.Account.OpenPositionCount.String()),
		"pending_order_count":            util.ParseFloat(res.Account.PendingOrderCount.String()),
		"pl":                             util.ParseFloat(res.Account.Pl),
		"unrealized_pl":                  util.ParseFloat(res.Account.UnrealizedPl),
		"nav":                            util.ParseFloat(res.Account.NAV),
		"margin_used":                    util.ParseFloat(res.Account.MarginUsed),
		"margin_available":               util.ParseFloat(res.Account.MarginAvailable),
		"position_value":                 util.ParseFloat(res.Account.PositionValue),
		"margin_closeout_unrealized_pl":  util.ParseFloat(res.Account.MarginCloseoutUnrealizedPL),
		"margin_closeout_nav":            util.ParseFloat(res.Account.MarginCloseoutNAV),
		"margin_closeout_margin_used":    util.ParseFloat(res.Account.MarginCloseoutMarginUsed),
		"margin_closeout_position_value": util.ParseFloat(res.Account.MarginCloseoutPositionValue),
		"margin_closeout_percent":        util.ParseFloat(res.Account.MarginCloseoutPercent),
	}
	return params
}

// APIのレスポンスをBusinessLogicのEntityに変換します。（単数銘柄）
func (a AccountManager) convertToEntityPosition(res *msg.PositionResponse) domain.Positions {
	units := util.ParseFloat(res.Position.Long.Units) + util.ParseFloat(res.Position.Short.Units)
	return domain.NewPosition(res.Position.Instrument, util.ParseFloat(res.Position.Pl), util.ParseFloat(res.Position.UnrealizedPL), util.ParseFloat(res.Position.MarginUsed), units)
}

// APIのレスポンスをBusinessLogicのEntityに変換します。（複数銘柄）
func (a AccountManager) convertToEntityPositions(res *msg.PositionsResponse) *[]domain.Positions {
	var positions []domain.Positions
	for _, p := range res.Positions {
		units := util.ParseFloat(p.Long.Units) + util.ParseFloat(p.Short.Units)
		position := domain.NewPosition(p.Instrument, util.ParseFloat(p.Pl), util.ParseFloat(p.UnrealizedPL), util.ParseFloat(p.MarginUsed), units)
		positions = append(positions, position)
	}
	return &positions
}
