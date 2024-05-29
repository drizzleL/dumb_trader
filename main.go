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
	"github.com/drizzleL/dumb_trader/chart"
	"github.com/drizzleL/dumb_trader/model"
	"github.com/drizzleL/dumb_trader/signal"
	"github.com/spf13/viper"
)

// 初始化币安api
func initApi() {
	viper.SetConfigFile("config.yaml")
	viper.ReadInConfig()
	apiKey := viper.GetString("apiKey")
	secretKey := viper.GetString("secretKey")
	binanceapi.Init(apiKey, secretKey)
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
}

// 打印图表
func PrintChart() {
	f, _ := os.Create("bar.html")
	klines := readProcess(model.Interval3M)
	klines2 := readProcess(model.Interval1M)
	chart.PrintEma(klines, f)
	chart.PrintEma(klines2, f)
}

func Mock() {
	// 用于开仓
	klines := readProcess(model.Interval3M)
	mainDict := map[int64]int{}
	for i := range klines.Original {
		mainDict[klines.Original[i].CloseTime] = i
	}

	// 小量级，用于平仓
	klines2 := readProcess(model.Interval1M)

	tr := &MockTrader{}
	for i := 0; i < len(klines2.CloseData); i++ {
		ts := klines2.Original[i].CloseTime
		if idx, ok := mainDict[ts]; ok {
			if consistent(klines, idx, 1) {
				tr.Long(klines.CloseData[idx].InexactFloat64(), klines.Original[idx].CloseTime)
			} else if consistent(klines, idx, -1) {
				tr.Short(klines.CloseData[idx].InexactFloat64(), klines.Original[idx].CloseTime)
			}
			continue
		}

		if klines2.Flag["break"][i] == 1 {
			tr.CloseShort(klines2.CloseData[i].InexactFloat64(), klines2.Original[i].CloseTime)
		} else if klines2.Flag["break"][i] == -1 {
			tr.CloseLong(klines2.CloseData[i].InexactFloat64(), klines2.Original[i].CloseTime)
		}

	}

	log.Println(tr.LastResult(klines.CloseData[len(klines.CloseData)-1].InexactFloat64()))
}

// 同时满足
func consistent(klines *model.Klines, i int, flag int) bool {
	if klines.Flag["break"][i] == flag && klines.Flag["cross"][i] == flag && klines.Flag["guppy"][i] == flag {
		return true
	}
	return false
}
