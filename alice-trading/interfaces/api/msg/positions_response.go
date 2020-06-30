package msg

type PositionsResponse struct {
	Positions         []Position `json:"positions"`
	LastTransactionID string     `json:"lastTransactionID"`
}

type PositionResponse struct {
	Position          Position `json:"position"`
	LastTransactionID string   `json:"lastTransactionID"`
}

type Position struct {
	Instrument              string       `json:"instrument"`
	Pl                      string       `json:"pl"`
	UnrealizedPL            string       `json:"unrealizedPL"`
	MarginUsed              string       `json:"marginUsed"`
	ResettablePL            string       `json:"resettablePL"`
	Financing               string       `json:"financing"`
	Commission              string       `json:"commission"`
	DividendAdjustment      string       `json:"dividendAdjustment"`
	GuaranteedExecutionFees string       `json:"guaranteedExecutionFees"`
	Long                    PositionSide `json:"long"`
	Short                   PositionSide `json:"short"`
}

type PositionSide struct {
	Units                   string   `json:"units"`
	AveragePrice            string   `json:"averagePrice"`
	TradeIDs                []string `json:"tradeIDs"`
	Pl                      string   `json:"pl"`
	UnrealizedPL            string   `json:"unrealizedPL"`
	ResettablePL            string   `json:"resettablePL"`
	Financing               string   `json:"financing"`
	DividendAdjustment      string   `json:"dividendAdjustment"`
	GuaranteedExecutionFees string   `json:"guaranteedExecutionFees"`
}
