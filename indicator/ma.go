package indicator

import "github.com/shopspring/decimal"

func Ma(klines []decimal.Decimal, n int) []decimal.Decimal {
	ret := make([]decimal.Decimal, len(klines))
	sum := decimal.New(0, 0)
	for i := 0; i < n; i++ { // 前n个不计算
		sum = sum.Add(klines[i])
	}
	div := decimal.NewFromInt(int64(n))
	for i := n; i < len(klines); i++ {
		sum = sum.Add(klines[i]).Sub(klines[i-n])
		ret[i] = sum.Div(div)
	}
	return ret
}
