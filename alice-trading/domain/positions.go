package domain

// ポジション
type Positions struct {
	ID            int
	Instrument    string
	Pl            float64
	UnrealizedPL  float64
	MarginUsed    float64
	Units         float64
	TradableShort bool
	TradableLong  bool
}

func NewPosition(instrument string, pl, unrealizedPL, marginUsed, units float64) Positions {
	// 逆建の注文は基本的に実施制限をかけるため、保有ポジションと同方向への取引のみ許可する
	tradableShort := true
	tradableLong := true
	switch {
	case units < 0:
		tradableLong = false
	case units >= 0:
		tradableShort = false
	}
	return Positions{
		Instrument:    instrument,
		Pl:            pl,
		UnrealizedPL:  unrealizedPL,
		MarginUsed:    marginUsed,
		Units:         units,
		TradableShort: tradableShort,
		TradableLong:  tradableLong,
	}

}
