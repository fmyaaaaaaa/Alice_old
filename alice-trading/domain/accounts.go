package domain

// アカウント
type Accounts struct {
	ID                          int
	MarginRate                  float64
	Balance                     float64
	OpenTradeCount              float64
	OpenPositionCount           float64
	PendingOrderCount           float64
	Pl                          float64
	UnrealizedPl                float64
	Nav                         float64
	MarginUsed                  float64
	MarginAvailable             float64
	PositionValue               float64
	MarginCloseoutUnrealizedPl  float64
	MarginCloseoutNav           float64
	MarginCloseoutMarginUsed    float64
	MarginCloseoutPositionValue float64
	MarginCloseoutPercent       float64
}
