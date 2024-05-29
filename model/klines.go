package model

import (
	"github.com/adshao/go-binance/v2/futures"
	"github.com/shopspring/decimal"
)

type Klines struct {
	Original  []*futures.Kline
	CloseData []decimal.Decimal
	Data      map[string][]decimal.Decimal
	Flag      map[string][]int
}

func NewKlines(original []*futures.Kline) *Klines {
	ret := &Klines{
		Original: original,
		Data:     make(map[string][]decimal.Decimal),
		Flag:     make(map[string][]int),
	}
	for _, line := range original {
		dec, _ := decimal.NewFromString(line.Close)
		ret.CloseData = append(ret.CloseData, dec)
	}
	return ret
}
