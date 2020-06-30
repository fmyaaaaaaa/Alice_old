package domain

// バックテストリザルト
type BackTestResults struct {
	ID          int
	TradeRule   string
	Instrument  string
	WinRate     float64
	LoseRate    float64
	MaxDrawDown float64
}
