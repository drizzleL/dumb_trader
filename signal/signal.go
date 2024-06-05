package signal

import (
	"github.com/drizzleL/dumb_trader/indicator"
	"github.com/drizzleL/dumb_trader/model"
)

// 添加信号
func Process(k *model.Klines) {
	k.Data["ma10"] = indicator.Ma(k.CloseData, 10)
	k.Data["ma20"] = indicator.Ma(k.CloseData, 20)
	k.Data["ma60"] = indicator.Ma(k.CloseData, 60)
	k.Data["ma120"] = indicator.Ma(k.CloseData, 120)

	k.Data["ema10"] = indicator.Ema(k.CloseData, 10)
	k.Data["ema20"] = indicator.Ema(k.CloseData, 20)
	k.Data["ema60"] = indicator.Ema(k.CloseData, 60)
	k.Data["ema120"] = indicator.Ema(k.CloseData, 120)

	k.Data["max5"] = indicator.Max(k.CloseData, 5)
	k.Data["min5"] = indicator.Min(k.CloseData, 5)

	k.Data["macd"] = indicator.Macd(k.CloseData)
	k.Data["sar"] = indicator.Sar(k.ProcessedData)

	addSignal(k)

	trimNoise(k)
}

func trimNoise(k *model.Klines) {
	k.Original = k.Original[120:]
	k.CloseData = k.CloseData[120:]
	k.ProcessedData = k.ProcessedData[120:]
	for key, val := range k.Flag {
		k.Flag[key] = val[120:]
	}
	for key, val := range k.Data {
		k.Data[key] = val[120:]
	}

}

func addSignal(k *model.Klines) {
	k.Flag["break"] = make([]int, len(k.CloseData))
	k.Flag["cross"] = make([]int, len(k.CloseData))
	k.Flag["guppy"] = make([]int, len(k.CloseData))
	k.Flag["conincr"] = make([]int, len(k.CloseData))
	k.Flag["condecr"] = make([]int, len(k.CloseData))
	k.Flag["max"] = make([]int, len(k.CloseData))
	k.Flag["min"] = make([]int, len(k.CloseData))
	for i := 0; i < len(k.CloseData); i++ {
		k.Flag["break"][i] = checkBreak(k, i)
		k.Flag["cross"][i] = checkCross(k, i)
		k.Flag["guppy"][i] = checkGuppy(k, i)
		k.Flag["conincr"][i] = checkConIncr(k, i)
		k.Flag["condecr"][i] = checkConDecr(k, i)
		k.Flag["max"][i] = checkMax(k, i)
		k.Flag["min"][i] = checkMin(k, i)
	}
}

// 连续增长
func checkConIncr(k *model.Klines, i int) int {
	close := k.CloseData[i]
	if i >= 1 {
		lastClose := k.CloseData[i-1]
		if close.LessThan(lastClose) {
			return 0
		}
	}
	if i >= 2 {
		lastClose := k.CloseData[i-2]
		if close.LessThan(lastClose) {
			return 0
		}
	}
	return 1
}

// 连续减少
func checkConDecr(k *model.Klines, i int) int {
	close := k.CloseData[i]
	if i >= 1 {
		lastClose := k.CloseData[i-1]
		if close.GreaterThan(lastClose) {
			return 0
		}
	}
	if i >= 2 {
		lastClose := k.CloseData[i-2]
		if close.GreaterThan(lastClose) {
			return 0
		}
	}
	return 1
}

// 收盘与短均线比较
func checkBreak(k *model.Klines, i int) int {
	close := k.CloseData[i]
	ma20 := k.Data["ma10"][i]
	ema20 := k.Data["ema10"][i]
	if close.GreaterThanOrEqual(ma20) && close.GreaterThanOrEqual(ema20) {
		return 1
	}
	if close.LessThanOrEqual(ma20) && close.LessThanOrEqual(ema20) {
		return -1
	}
	return 0
}

func checkMax(k *model.Klines, i int) int {
	close := k.CloseData[i]
	maxVal := k.Data["max5"][i]
	if close.GreaterThanOrEqual(maxVal) {
		return 1
	}
	return 0
}

func checkMin(k *model.Klines, i int) int {
	close := k.CloseData[i]
	maxVal := k.Data["min5"][i]
	if close.LessThanOrEqual(maxVal) {
		return 1
	}
	return 0
}

// 短均线与中均线比较
func checkCross(k *model.Klines, i int) int {
	ma20 := k.Data["ma20"][i]
	ma60 := k.Data["ma60"][i]
	ema20 := k.Data["ema20"][i]
	ema60 := k.Data["ema60"][i]
	if ma20.GreaterThanOrEqual(ma60) && ema20.GreaterThanOrEqual(ema60) {
		return 1
	}
	if ma20.LessThanOrEqual(ma60) && ema20.LessThanOrEqual(ema60) {
		return -1
	}
	return 0
}

// 短均线与中均线，长均线比较
func checkGuppy(k *model.Klines, i int) int {
	// ma20 := k.Data["ma20"][i]
	// ma60 := k.Data["ma60"][i]
	// ma120 := k.Data["ma120"][i]
	ema20 := k.Data["ema20"][i]
	ema60 := k.Data["ema60"][i]
	ema120 := k.Data["ema120"][i]
	if ema20.GreaterThanOrEqual(ema60) && ema20.GreaterThanOrEqual(ema120) {
		return 1
	}
	if ema20.LessThanOrEqual(ema60) && ema20.LessThanOrEqual(ema120) {
		return -1
	}
	return 0
}
