package chart

import (
	"github.com/go-echarts/go-echarts/v2/actions"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func generateKline(chartData [][]interface{}) *charts.Kline {

	kline := charts.NewKLine()
	kline.AddDataset(opts.Dataset{Source: chartData})
	kline.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Theme: "dark",
		}),
		charts.WithGridOpts(
			opts.Grid{Bottom: "210", Left: "50", Right: "10"},
			opts.Grid{Height: "80", Bottom: "210", Left: "50", Right: "10"},
			opts.Grid{Height: "80", Bottom: "80", Left: "50", Right: "10"},
			opts.Grid{Height: "80", Bottom: "10", Left: "50", Right: "10"}),
		charts.WithTitleOpts(opts.Title{
			Title: "Kline-example",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Type:      "category",
			Scale:     true,
			GridIndex: 0,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Type:      "value",
			Scale:     true,
			GridIndex: 0,
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "inside",
			XAxisIndex: []int{0, 1, 2},
			Start:      90,
			End:        100,
		},
			opts.DataZoom{
				Type:       "slider",
				XAxisIndex: []int{0, 1, 2},
				Start:      98,
				End:        100,
			}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show: true, Trigger: "axis",
			AxisPointer: &opts.AxisPointer{Type: "line"},
		}),
		charts.WithToolboxOpts(opts.Toolbox{
			Show: true,
			Feature: &opts.ToolBoxFeature{
				DataZoom: &opts.ToolBoxFeatureDataZoom{YAxisIndex: false},
				Brush:    &opts.ToolBoxFeatureBrush{Type: []string{"lineX", "clear"}},
			},
		}),
		charts.WithBrush(opts.Brush{
			XAxisIndex: "all",
			Brushlink:  "all",
			OutOfBrush: &opts.BrushOutOfBrush{ColorAlpha: 0.1},
		}),
	)
	kline.SetDispatchActions(
		charts.WithAreas(actions.Areas{
			BrushType:  "lineX",
			CoordRange: []string{"2020-08-02 18:00:00", "2020-08-02 10:00:00"},
			XAxisIndex: 0,
		}),
		charts.WithType("brush"),
	)
	kline.AddSeries("candlestick", nil, charts.WithEncodeOpts(opts.Encode{X: "date", Y: [4]string{"open", "close", "low", "high"}}))
	kline.ExtendXAxis(opts.XAxis{Type: "category", SplitNumber: 20, GridIndex: 1,
		AxisTick:  &opts.AxisTick{Show: false},
		AxisLabel: &opts.AxisLabel{Show: false},
	})
	kline.ExtendYAxis(opts.YAxis{
		Scale: true, GridIndex: 1, SplitNumber: 2,
		AxisLabel: &opts.AxisLabel{Show: false},
		AxisLine:  &opts.AxisLine{Show: false},
		SplitLine: &opts.SplitLine{Show: false},
	})

	kline.ExtendXAxis(opts.XAxis{Type: "category", SplitNumber: 20, GridIndex: 2,
		AxisTick:  &opts.AxisTick{Show: false},
		AxisLabel: &opts.AxisLabel{Show: false},
	})
	kline.ExtendYAxis(opts.YAxis{
		Scale: true, GridIndex: 2, SplitNumber: 2,
		AxisLabel: &opts.AxisLabel{Show: true},
		AxisLine:  &opts.AxisLine{Show: true},
		SplitLine: &opts.SplitLine{Show: true},
	})

	volumeBarChart := charts.NewBar()
	volumeBarChart.AddSeries("Volume", nil,
		charts.WithItemStyleOpts(opts.ItemStyle{Color: "#7fbe9e"}),
		charts.WithBarChartOpts(opts.BarChart{Type: "bar", XAxisIndex: 1, YAxisIndex: 1}),
		charts.WithEncodeOpts(opts.Encode{X: "date", Y: "volume"}))

	ema10LineChart := charts.NewLine()
	ema10LineChart.SetGlobalOptions(charts.WithXAxisOpts(opts.XAxis{SplitNumber: 20, GridIndex: 0}), charts.WithYAxisOpts(opts.YAxis{Scale: true, GridIndex: 0}))
	ema10LineChart.AddSeries("EMA10", nil,
		charts.WithEncodeOpts(opts.Encode{X: "date", Y: "ema10"}),
		charts.WithLineStyleOpts(opts.LineStyle{Color: "rgba(69, 140, 255, 0.5)"}),
		charts.WithItemStyleOpts(opts.ItemStyle{Opacity: 0.01}),
		charts.WithLineChartOpts(opts.LineChart{XAxisIndex: 0, YAxisIndex: 0}))

	rsChart := charts.NewLine()
	rsChart.AddSeries("RS", nil,
		charts.WithEncodeOpts(opts.Encode{X: "date", Y: "rs"}),
		charts.WithLineChartOpts(opts.LineChart{XAxisIndex: 2, YAxisIndex: 2}),
		charts.WithLineStyleOpts(opts.LineStyle{Color: "#0f0"}),
		charts.WithItemStyleOpts(opts.ItemStyle{Opacity: 0.01}),
	)
	kline.Overlap(rsChart)

	rsiChart := charts.NewLine()
	rsiChart.AddSeries("RSI", nil,
		charts.WithEncodeOpts(opts.Encode{X: "date", Y: "rsi"}),
		charts.WithLineChartOpts(opts.LineChart{XAxisIndex: 2, YAxisIndex: 2}),
		charts.WithLineStyleOpts(opts.LineStyle{Color: "rgba(169, 84, 255, 0.5)"}),
		charts.WithItemStyleOpts(opts.ItemStyle{Opacity: 0.01}))
	kline.Overlap(rsiChart)

	kline.Overlap(volumeBarChart)
	kline.Overlap(ema10LineChart)
	return kline
}
