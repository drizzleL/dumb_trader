package indicator

import (
	"github.com/shopspring/decimal"
)

func Rsi(klines []decimal.Decimal) []decimal.Decimal {
	uu := make([]decimal.Decimal, len(klines))
	dd := make([]decimal.Decimal, len(klines))
	for i := 1; i < len(klines); i++ {
		if klines[i].GreaterThan(klines[i-1]) {
			uu[i] = klines[i].Sub(klines[i-1])
		} else {
			dd[i] = klines[i-1].Sub(klines[i])
		}
	}
	uuEma := Ema(uu, 6)
	ddEma := Ema(dd, 6)
	ret := make([]decimal.Decimal, len(klines))
	one := decimal.NewFromInt(1)
	hundred := decimal.NewFromInt(100)
	for i := 6; i < len(klines); i++ {
		rs := uuEma[i].Div(ddEma[i])
		ret[i] = rs.Div(rs.Add(one)).Mul(hundred)
	}
	return ret
}
