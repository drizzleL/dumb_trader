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
	ep := klines[0].High
	flag := sarFlagIncr
	ret[0] = klines[0].Low
	if klines[0].Open.GreaterThan(klines[0].Close) {
		flag = sarFlagDecr
		ret[0] = klines[0].High
		ep = klines[0].Low
	}
	af := incrFactor
	for i := 1; i < len(klines); i++ {
		ret[i] = ret[i-1].Sub(af.Mul(ret[i-1].Sub(ep)))
		switch flag {
		case sarFlagIncr:
			if klines[i].High.GreaterThan(ep) {
				af = addFactor(af)
				ep = klines[i].High
			}
			if klines[i].Low.LessThan(ret[i]) { // reverse
				ret[i] = ep
				ep = klines[i].Low
				flag = sarFlagDecr
				af = incrFactor
			}
		case sarFlagDecr:
			if klines[i].Low.LessThan(ep) {
				af = addFactor(af)
				ep = klines[i].Low
			}
			if klines[i].High.GreaterThan(ret[i]) { // reverse
				ret[i] = ep
				ep = klines[i].High
				flag = sarFlagIncr
				af = incrFactor
			}
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
