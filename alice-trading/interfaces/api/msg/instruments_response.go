package msg

import "encoding/json"

type InstrumentsResponse struct {
	Instruments []Instrument
}

type Instrument struct {
	DisplayName                 string      `json:"displayName"`
	DisplayPrecision            json.Number `json:"displayPrecision"`
	MarginRate                  string      `json:"marginRate"`
	MaximumOrderUnits           string      `json:"maximumOrderUnits"`
	MaximumPositionSize         string      `json:"maximumPositionSize"`
	MaximumTrailingStopDistance string      `json:"maximumTrailingStopDistance"`
	MinimumTradeSize            string      `json:"minimumTradeSize"`
	MinimumTrailingStopDistance string      `json:"minimumTrailingStopDistance"`
	Name                        string      `json:"name"`
	PipLocation                 json.Number `json:"pipLocation"`
	TradeUnitsPrecision         json.Number `json:"tradeUnitsPrecision"`
	Type                        string      `json:"type"`
}
