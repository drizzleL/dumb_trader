package indicator

import (
	"github.com/drizzleL/dumb_trader/model"
	"github.com/shopspring/decimal"
)

type sarFlag int

const (
	sarFlagIncr = iota
	sarFlagDecr
)

func Sar(klines []*model.Data) []decimal.Decimal {
	ret := make([]decimal.Decimal, len(klines))
	// 找到第一个极值
	low, high := klines[0].Low, klines[0].High
	flag := sarFlagIncr
	ret[0] = klines[0].Low
	if klines[0].Open.GreaterThan(klines[0].Close) {
		flag = sarFlagDecr
		ret[0] = klines[0].High
	}
	af := incrFactor
	for i := 1; i < len(klines); i++ {
		switch flag {
		case sarFlagIncr:
			if klines[i].Low.LessThan(ret[i-1]) { // reverse
				low, high = klines[i].Low, decimal.Max(high, klines[i].High)
				ret[i] = high
				flag = sarFlagDecr
				af = incrFactor
				continue
			}
			if klines[i].High.GreaterThan(high) {
				af = addFactor(af)
			}
			low, high = klines[i].Low, decimal.Max(high, klines[i].High)
			ret[i] = ret[i-1].Sub(af.Mul(ret[i-1].Sub(high)))
		case sarFlagDecr:
			if klines[i].High.GreaterThan(ret[i-1]) { // reverse
				low, high = decimal.Min(low, klines[i].Low), klines[i].High
				ret[i] = low
				flag = sarFlagIncr
				af = incrFactor
				continue
			}
			if klines[i].Low.LessThan(low) {
				af = addFactor(af)
			}
			low, high = decimal.Min(low, klines[i].Low), klines[i].High
			ret[i] = ret[i-1].Sub(af.Mul(ret[i-1].Sub(low)))
		}
	}
	return ret
}

var (
	maxFactor  = decimal.NewFromFloat(0.2)
	incrFactor = decimal.NewFromFloat(0.02)
)

func addFactor(af decimal.Decimal) decimal.Decimal {
	if af.Equal(maxFactor) {
		return af
	}
	return af.Add(incrFactor)
}
