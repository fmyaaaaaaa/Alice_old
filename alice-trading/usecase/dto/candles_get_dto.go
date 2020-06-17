package dto

import "github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"

// 足データ取得DTO
type CandlesGetDto struct {
	InstrumentName string           //通貨名
	Count          int              //本数
	Granularity    enum.Granularity //足種
}

func NewCandlesGetDto(instrumentName string, count int, granularity enum.Granularity) CandlesGetDto {
	return CandlesGetDto{
		InstrumentName: instrumentName,
		Count:          count,
		Granularity:    granularity,
	}
}
