package instruments

type Instruments struct {
	ID                         int
	Name                       string
	PipLocation                float32
	MaximumOrderUnits          float32
	MinimumTradeSize           float32
	MinimumTradingStopDistance float32
}
