package binanceapi

import (
	"context"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/drizzleL/dumb_trader/model"
)

type CollectReq struct {
	Symbol    string
	StartTime time.Time
	EndTime   time.Time
	Interval  *model.Interval
}

func Collect(req CollectReq) []*futures.Kline {
	var ret []*futures.Kline
	start := req.StartTime
	for start.Before(req.EndTime) {
		end := start.Add(500 * req.Interval.Duration)
		srv := futuresClient.NewKlinesService()
		srv.Symbol(req.Symbol)
		srv.Interval(req.Interval.Name)

		srv.StartTime(start.UnixMilli())
		srv.EndTime(end.UnixMilli() - 1)

		resp, err := srv.Do(context.Background())
		if err != nil {
			panic(err)
		}
		ret = append(ret, resp...)
		start = end
	}
	return ret
}
