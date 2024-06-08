package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/drizzleL/dumb_trader/binanceapi"
	"github.com/drizzleL/dumb_trader/bot"
	"github.com/drizzleL/dumb_trader/chart"
	"github.com/drizzleL/dumb_trader/model"
	"github.com/drizzleL/dumb_trader/signal"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile("config.yaml")
	viper.ReadInConfig()
}

// 初始化币安api
func initApi() {
	apiKey := viper.GetString("apiKey")
	secretKey := viper.GetString("secretKey")
	binanceapi.Init(apiKey, secretKey)
}

func initBot() {
	bot.Init(viper.GetString("botToken"))
}

// 将范围内k线数据写入文件
func writeFile(intv model.Interval) {
	data := binanceapi.Collect(binanceapi.CollectReq{
		Symbol:    "DOGEUSDT",
		Interval:  intv,
		StartTime: time.Unix(1716307200, 0),
		EndTime:   time.Unix(1716912000, 0),
	})
	b, _ := json.Marshal(data)
	fname := fmt.Sprintf("%s.json", intv.Name)
	f, _ := os.Create(fname)
	f.Write(b)
	f.Close()
}

func readFile(intv model.Interval) []*futures.Kline {
	fname := fmt.Sprintf("%s.json", intv.Name)
	f, _ := os.Open(fname)
	b, _ := io.ReadAll(f)
	var dest []*futures.Kline
	json.Unmarshal(b, &dest)
	return dest
}

// 通过api获取json数据
func LoadData() {
	initApi()
	writeFile(model.Interval1M)
	writeFile(model.Interval3M)
	writeFile(model.Interval5M)
	writeFile(model.Interval15M)
}

func readProcess(intv model.Interval) *model.Klines {
	original := readFile(intv)
	klines := model.NewKlines(original)
	signal.Process(klines)
	return klines
}

func main() {
	Mock()
	PrintChart()
}

// 打印图表
func PrintChart() {
	f, _ := os.Create("bar.html")
	klines := readProcess(model.Interval15M)
	// klines2 := readProcess(model.Interval3M)

	chart.PrintKline(klines, f)

	// chart.PrintEma(klines, f)
	// chart.PrintEma(klines2, f)
}

func Mock() {
	// 用于开仓
	klines := readProcess(model.Interval15M)
	klines2 := readProcess(model.Interval3M)

	group := NewKlineGroup(klines, klines2)

	str := &BreakStrategy{}

	tr := &MockTrader{}
	for i := 0; i < len(klines2.CloseData); i++ {
		ts := klines2.Original[i].CloseTime
		opType, price := str.CheckDeal(group, ts)
		switch opType {
		case OpTypeLong:
			tr.Long(price, ts)
		case OpTypeShort:
			tr.Short(price, ts)
		case OpTypeCloseLong:
			tr.CloseLong(price, ts)
		case OpTypeCloseShort:
			tr.CloseShort(price, ts)
		}
	}

	log.Println(tr.LastResult(klines.CloseData[len(klines.CloseData)-1].InexactFloat64()))
}

// 同时满足
func consistent(klines *model.Klines, i int, flag int) bool {
	if klines.Flag["guppy"][i] == flag {
		return true
	}
	return false
}

type KlineGroup struct {
	MainKlines *model.Klines
	SideKlines *model.Klines
	MainDict   map[int64]int
	SideDict   map[int64]int
}

func NewKlineGroup(main, side *model.Klines) *KlineGroup {
	mainDict := map[int64]int{}
	for i, v := range main.Original {
		mainDict[v.CloseTime] = i
	}
	sideDict := map[int64]int{}
	for i, v := range side.Original {
		sideDict[v.CloseTime] = i
	}
	return &KlineGroup{
		MainKlines: main,
		SideKlines: side,
		MainDict:   mainDict,
		SideDict:   sideDict,
	}

}

type Strategy interface {
	CheckDeal(klines *KlineGroup, ts int64) OpType
}

type OpType int

const (
	OpTypeDefault OpType = iota
	OpTypeLong
	OpTypeShort
	OpTypeCloseLong
	OpTypeCloseShort
)

type BreakStrategy struct{}

func (b BreakStrategy) CheckDeal(group *KlineGroup, ts int64) (OpType, float64) {
	// sideIdx := group.SideDict[ts]
	mainIdx, ok := group.MainDict[ts]
	if ok && group.MainKlines.CloseData[mainIdx].LessThan(group.MainKlines.Data["sar"][mainIdx]) {
		return OpTypeCloseLong, group.MainKlines.CloseData[mainIdx].InexactFloat64()
	}
	if ok && group.MainKlines.Flag["guppy"][mainIdx] == 1 {
		return OpTypeLong, group.MainKlines.CloseData[mainIdx].InexactFloat64()
	}
	// if ok && group.MainKlines.Flag["guppy"][mainIdx] == -1 {
	// 	return OpTypeShort, group.MainKlines.CloseData[mainIdx].InexactFloat64()
	// }

	// if ok && group.SideKlines.Flag["break"][sideIdx] == 1 {
	// return OpTypeCloseShort, group.SideKlines.CloseData[sideIdx].InexactFloat64()
	// }

	// if ok && group.MainKlines.CloseData[mainIdx].GreaterThan(group.MainKlines.Data["sar"][mainIdx]) {
	// 	return OpTypeCloseShort, group.MainKlines.CloseData[mainIdx].InexactFloat64()
	// }
	return OpTypeDefault, 0
}
