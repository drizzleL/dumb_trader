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

func NewKlines(original []*futures.Kline, subData []*futures.Kline) *Klines {
	ret := &Klines{
		Original: original,
		Data:     make(map[string][]decimal.Decimal),
		Flag:     make(map[string][]int),
	}
	var subJ int
	for _, line := range original {
		dec, _ := decimal.NewFromString(line.Close)
		ret.CloseData = append(ret.CloseData, dec)

		d := klineToData(line)
		for subJ < len(subData) && subData[subJ].CloseTime < line.OpenTime {
			subJ += 1
		}
		for subJ < len(subData) && subData[subJ].OpenTime >= line.OpenTime && subData[subJ].CloseTime <= line.CloseTime {
			d.ChildData = append(d.ChildData, klineToData(subData[subJ]))
			subJ += 1
		}
		ret.ProcessedData = append(ret.ProcessedData, d)
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
		Volume:    strToDec(line.Volume),
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
	OpenTime  int64
	CloseTime int64
	Volume    decimal.Decimal
	ChildData []*Data
}
