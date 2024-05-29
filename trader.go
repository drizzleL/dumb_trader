package main

import (
	"log"
	"time"
)

type MockTrader struct {
	Pay      float64
	Receive  float64
	HoldFlag int // 看多1，空-1
}

func (d *MockTrader) LastResult(price float64) float64 {
	if d.HoldFlag == 0 {
		return d.Receive - d.Pay
	}
	switch d.HoldFlag {
	case 1:
		d.Receive += price
	case -1:
		d.Pay += price
	}

	return d.Receive - d.Pay
}

var (
	opPrintLong       = "做多"
	opPrintShort      = "做空"
	opPrintCloseLong  = "平多"
	opPrintCloseShort = "平空"
)

func printWithTime(op string, price float64, t int64) {
	ts := time.UnixMilli(t).Format("2006-01-02T15:04:05Z")
	log.Printf("%s %s at %f", ts, op, price)
}

func (d *MockTrader) Long(price float64, t int64) {
	if d.HoldFlag == 0 {
		d.HoldFlag = 1
		d.Pay += price
		printWithTime(opPrintLong, price, t)
		return
	}
	if d.HoldFlag == 1 {
		return
	}
	return
}

func (d *MockTrader) Short(price float64, t int64) {
	if d.HoldFlag == 0 {
		d.HoldFlag = -1
		d.Receive += price
		printWithTime(opPrintShort, price, t)
		return
	}
	if d.HoldFlag == -1 {
		return
	}
	return
}

func (d *MockTrader) CloseLong(price float64, t int64) {
	if d.HoldFlag != 1 {
		return
	}
	d.Receive += price
	d.HoldFlag = 0
	printWithTime(opPrintCloseLong, price, t)
	println()
}
func (d *MockTrader) CloseShort(price float64, t int64) {
	if d.HoldFlag != -1 {
		return
	}
	d.Pay += price
	d.HoldFlag = 0
	printWithTime(opPrintCloseShort, price, t)
	println()
}
