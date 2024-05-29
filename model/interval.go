package model

import "time"

var (
	Interval1M = Interval{
		Duration: time.Minute * 1,
		Name:     "1m",
	}
	Interval3M = Interval{
		Duration: time.Minute * 3,
		Name:     "3m",
	}
	Interval5M = Interval{
		Duration: time.Minute * 5,
		Name:     "5m",
	}
	Interval15M = Interval{
		Duration: time.Minute * 15,
		Name:     "15m",
	}
)

type Interval struct {
	Duration time.Duration
	Name     string
}
