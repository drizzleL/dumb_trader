package model

import (
	"github.com/adshao/go-binance/v2/futures"
	"github.com/shopspring/decimal"
)

type Klines struct {
	Original      []*futures.Kline
	ProcessedData []*Data
	CloseData     []decimal.Decimal
	Data          map[string][]decimal.Decimal
	Flag          map[string][]int
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

		ret.ProcessedData = append(ret.ProcessedData, klineToData(line))
	}
	return ret
}

func klineToData(line *futures.Kline) *Data {
	return &Data{
		CloseTime: line.CloseTime,
		Open:      strToDec(line.Open),
		Close:     strToDec(line.Close),
		High:      strToDec(line.High),
		Low:       strToDec(line.Low),
	}
}

func strToDec(s string) decimal.Decimal {
	d, _ := decimal.NewFromString(s)
	return d
}

type Data struct {
	High      decimal.Decimal
	Low       decimal.Decimal
	Open      decimal.Decimal
	Close     decimal.Decimal
	CloseTime int64
}
