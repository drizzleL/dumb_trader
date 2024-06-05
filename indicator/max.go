package indicator

import "github.com/shopspring/decimal"

func Max(klines []decimal.Decimal, k int) []decimal.Decimal {
	var stack []decimal.Decimal
	ret := make([]decimal.Decimal, len(klines))
	for i := 0; i < k; i++ {
		for len(stack) != 0 && stack[len(stack)-1].LessThan(klines[i]) {
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, klines[i])
		ret[i] = stack[0]
	}
	for i := k; i < len(klines); i++ {
		if klines[i-k].Equal(stack[0]) {
			stack = stack[1:]
		}
		for len(stack) != 0 && stack[len(stack)-1].LessThan(klines[i]) {
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, klines[i])
		ret[i] = stack[0]
	}
	return ret
}

func Min(klines []decimal.Decimal, k int) []decimal.Decimal {
	var stack []decimal.Decimal
	ret := make([]decimal.Decimal, len(klines))
	for i := 0; i < k; i++ {
		for len(stack) != 0 && stack[len(stack)-1].GreaterThan(klines[i]) {
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, klines[i])
		ret[i] = stack[0]
	}
	for i := k; i < len(klines); i++ {
		if klines[i-k].Equal(stack[0]) {
			stack = stack[1:]
		}
		for len(stack) != 0 && stack[len(stack)-1].GreaterThan(klines[i]) {
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, klines[i])
		ret[i] = stack[0]
	}
	return ret
}
