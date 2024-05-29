package indicator

import "github.com/shopspring/decimal"

func Ema(klines []decimal.Decimal, n int) []decimal.Decimal {
	ret := make([]decimal.Decimal, len(klines))
	sum := decimal.New(0, 0)
	for i := 0; i < n; i++ { // 前n个不计算
		sum = sum.Add(klines[i])
	}
	ret[n-1] = sum.Div(decimal.NewFromInt(int64(n)))
	div := decimal.NewFromInt(int64(n) + 1)
	twoDecimal := decimal.NewFromInt(2)
	dec2 := decimal.NewFromInt(int64(n - 1))
	for i := n; i < len(klines); i++ {
		ret[i] = klines[i].Mul(twoDecimal).Add(ret[i-1].Mul(dec2)).Div(div)
	}
	return ret
}
