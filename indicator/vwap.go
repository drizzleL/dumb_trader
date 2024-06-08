package indicator

import (
	"github.com/drizzleL/dumb_trader/model"
	"github.com/shopspring/decimal"
)

func Vwap(klines []*model.Data) []decimal.Decimal {
	ret := make([]decimal.Decimal, len(klines))
	var sum, volSum decimal.Decimal
	threeDec := decimal.NewFromInt(3)
	for i, line := range klines {
		midPrice := line.Close.Add(line.High).Add(line.Low).Div(threeDec)
		sum = sum.Add(midPrice.Mul(line.Volume))
		volSum = volSum.Add(line.Volume)
		ret[i] = sum.Div(volSum)
	}
	return ret
}
