package signal

import (
	"github.com/drizzleL/dumb_trader/model"
	"github.com/shopspring/decimal"
)

var threeDec = decimal.NewFromInt(3)

func checkGather(k *model.Klines, i int) decimal.Decimal {
	return k.Data["ma20"][i].Add(k.Data["ma60"][i]).Add(k.Data["ma120"][i]).Div(threeDec)
}
