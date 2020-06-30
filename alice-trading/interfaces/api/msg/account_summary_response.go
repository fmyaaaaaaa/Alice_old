package msg

import "encoding/json"

type AccountSummaryResponse struct {
	Account AccountSummary `json:"account"`
}

type AccountSummary struct {
	ID                          string      `json:"id"`
	MarginRate                  string      `json:"marginRate"`
	Balance                     string      `json:"balance"`
	OpenTradeCount              json.Number `json:"openTradeCount"`
	OpenPositionCount           json.Number `json:"openPositionCount"`
	PendingOrderCount           json.Number `json:"pendingOrderCount"`
	Pl                          string      `json:"pl"`
	UnrealizedPl                string      `json:"unrealizedPL"`
	NAV                         string      `json:"NAV"`
	MarginUsed                  string      `json:"marginUsed"`
	MarginAvailable             string      `json:"marginAvailable"`
	PositionValue               string      `json:"positionValue"`
	MarginCloseoutUnrealizedPL  string      `json:"marginCloseoutUnrealizedPL"`
	MarginCloseoutNAV           string      `json:"marginCloseoutNAV"`
	MarginCloseoutMarginUsed    string      `json:"marginCloseoutMarginUsed"`
	MarginCloseoutPositionValue string      `json:"marginCloseoutPositionValue"`
	MarginCloseoutPercent       string      `json:"marginCloseoutPercent"`
	WithdrawalLimit             string      `json:"withdrawalLimit"`
}
