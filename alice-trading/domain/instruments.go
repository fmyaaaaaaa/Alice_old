package domain

type Instruments struct {
	ID                         int
	Instrument                 string
	PipLocation                float64
	MinimumTradeSize           float64
	MaximumTradingStopDistance float64
	MinimumTradingStopDistance float64
	MaximumPositionSize        float64
	MaximumOrderUnits          float64
	MarginRate                 float64
	EvaluationInstrument       string
}

func (i Instruments) IsJpyEvaluation() bool {
	if i.EvaluationInstrument == "JPY" {
		return true
	}
	return false
}
