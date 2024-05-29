package indicator

import "github.com/shopspring/decimal"

func Macd(klines []decimal.Decimal) []decimal.Decimal {
	ema1 := Ema(klines, 12)
	ema2 := Ema(klines, 26)
	ret := make([]decimal.Decimal, len(klines))
	sum := decimal.New(0, 0)
	for i := 27; i < 36; i++ {
		dif := ema1[i].Sub(ema2[i])
		sum = sum.Add(dif)
	}
	ret[35] = sum.Div(decimal.NewFromInt(9))
	div := decimal.NewFromInt(10)
	twoDecimal := decimal.NewFromInt(2)
	dec2 := decimal.NewFromInt(8)
	for i := 36; i < len(klines); i++ {
		dif := ema1[i].Sub(ema2[i])
		last := ret[i-1]
		ret[i] = dif.Mul(twoDecimal).Add(last.Mul(dec2)).Div(div)
	}
	return ret
}
