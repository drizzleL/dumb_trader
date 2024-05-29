package trade

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/drizzleL/dumb_trader/binanceapi"
	"github.com/shopspring/decimal"
)

var callbackRate = "0.3"

var money int
var opType int
var strategy string

var s5 = &Strategy{
	ProfitPercent:  15,
	ProfitPercent2: 40,
	CallbackRate:   10,
	LossPercent:    20,
}

// quick用于快速波动
var s1 = &Strategy{
	ProfitPercent:  25,
	ProfitPercent2: 40,
	CallbackRate:   10,
	LossPercent:    20,
}

// slow用于平稳曲线，5%
var s2 = &Strategy{
	ProfitPercent: 20,
	CallbackRate:  15,
	LossPercent:   50,
}

// wait用于磨底, 30%
var s3 = &Strategy{
	ProfitPercent: 30,
	CallbackRate:  10,
	LossPercent:   15,
}

// chase追涨杀跌
var s4 = &Strategy{
	ProfitPercent: 5,
	CallbackRate:  5,
}

const tradeMargin = 50

var apiSwitch bool = true

var (
	symbolDoge = "DOGEUSDT"
	symbolEna  = "ENAUSDT"
)

func quick(strategy string) {
	var s *Strategy
	switch strategy {
	case "quick":
		s = s1
	case "slow":
		s = s2
	case "wait":
		s = s3
	case "chase":
		s = s4
	case "fast":
		s = s5
	}
	switch opType {
	case 1:
		buyFutures(context.Background(), &BuyRequest{
			Symbol: symbolDoge,
			Money:  money,
		}, s)
	case 2:
		sellFutures(context.Background(), &BuyRequest{
			Symbol: symbolDoge,
			Money:  money,
		}, s)
	case 3:
		buyFutures(context.Background(), &BuyRequest{
			Symbol: symbolEna,
			Money:  money,
		}, s)
	case 4:
		sellFutures(context.Background(), &BuyRequest{
			Symbol: symbolEna,
			Money:  money,
		}, s)
	}
}

type BuyRequest struct {
	Symbol  string
	Money   int // usdt as unit
	BuySell int
}

type Strategy struct {
	ProfitPercent  int
	LossPercent    int
	CallbackRate   int
	ProfitPercent2 int
}

type Offer struct {
	Symbol       string
	Price        string
	Amount       string
	CallbackRate int
}

func sellFutures(ctx context.Context, req *BuyRequest, s *Strategy) {
	price := binanceapi.ShowPrice(ctx, req.Symbol)
	dec, err := decimal.NewFromString(price)
	if err != nil {
		panic(err)
	}
	priceDec, err := decimal.NewFromString(price)
	if err != nil {
		panic(err)
	}
	amount := decimal.NewFromInt(int64(req.Money * tradeMargin)).Div(priceDec).Round(0).String()
	Buy(ctx, &Offer{
		Symbol: req.Symbol,
		Amount: amount,
	}, TradeSell)
	if s != nil {
		if s.ProfitPercent != 0 {
			diff := getDiff(dec, s.ProfitPercent).Round(5)
			profit := dec.Sub(diff).Round(5).String()
			log.Printf("[BUY] profit_percent: %d%%, at %s", s.ProfitPercent, profit)
			CloseTradeTrail(ctx, &Offer{
				Price:        profit,
				Symbol:       req.Symbol,
				Amount:       amount,
				CallbackRate: s.CallbackRate,
			}, TradeBuy)
		}
		if s.LossPercent != 0 {
			diff := getDiff(dec, s.LossPercent).Round(5)
			loss := dec.Add(diff).Round(5).String()
			log.Printf("[BUY] loss_percent: %d%%, at %s", s.LossPercent, loss)
			CloseTrade(ctx, &Offer{
				Price:  loss,
				Symbol: req.Symbol,
			}, TradeBuy, CloseStop)
		}
	}
}

// long doge 1000usdt, with strategy_1
func buyFutures(ctx context.Context, req *BuyRequest, s *Strategy) {
	price := binanceapi.ShowPrice(ctx, req.Symbol)
	priceDec, err := decimal.NewFromString(price)
	if err != nil {
		panic(err)
	}
	amount := decimal.NewFromInt(int64(req.Money * tradeMargin)).Div(priceDec).Round(0).String()
	Buy(ctx, &Offer{
		Symbol: req.Symbol,
		Amount: amount,
	}, TradeBuy)
	log.Printf("[BUY] %s, current price: %s, cost: $%d", req.Symbol, price, req.Money)
	if s != nil {
		dec, err := decimal.NewFromString(price)
		if err != nil {
			panic(err)
		}
		if s.ProfitPercent != 0 {
			diff := getDiff(dec, s.ProfitPercent).Round(5)
			profit := dec.Add(diff).Round(5).String()
			log.Printf("[BUY] profit_percent: %d%%, at %s", s.ProfitPercent, profit)
			CloseTradeTrail(ctx, &Offer{
				Price:        profit,
				Symbol:       req.Symbol,
				Amount:       amount,
				CallbackRate: s.CallbackRate,
			}, TradeSell)
		}
		if s.LossPercent != 0 {
			diff := getDiff(dec, s.LossPercent).Round(5)
			loss := dec.Sub(diff).Round(5).String()
			log.Printf("[BUY] loss_percent: %d%%, at %s", s.LossPercent, loss)
			CloseTrade(ctx, &Offer{
				Price:  loss,
				Symbol: req.Symbol,
			}, TradeSell, CloseStop)
		}
	}
}

func getDiff(price decimal.Decimal, percent int) decimal.Decimal {
	v := decimal.NewFromInt(int64(percent))
	ratio := v.Div(decPercent).Div(decimal.NewFromInt(tradeMargin))
	return price.Mul(ratio)
}

var decPercent = decimal.NewFromInt(100)

type CloseType int

const (
	CloseProfit CloseType = iota
	CloseStop
)

type TradeType int

const (
	TradeSell TradeType = iota
	TradeBuy
)

var futuresClient *futures.Client

func Buy(ctx context.Context, offer *Offer, tradeType TradeType) {
	log.Printf("buy symbol: %s, amount: %s", offer.Symbol, offer.Amount)
	if !apiSwitch {
		return
	}
	srv := futuresClient.NewCreateOrderService()
	srv.Symbol(offer.Symbol)
	srv.Type(futures.OrderTypeMarket).Quantity(offer.Amount)
	switch tradeType {
	case TradeBuy:
		srv.Side(futures.SideTypeBuy)
		srv.PositionSide(futures.PositionSideTypeLong)
	case TradeSell:
		srv.Side(futures.SideTypeSell)
		srv.PositionSide(futures.PositionSideTypeShort)
	}
	res, err := srv.Do(ctx)
	if err != nil {
		panic(err)
	}
	printVal(res)
}

func Shut(ctx context.Context, symbol string) {
	if !apiSwitch {
		return
	}
	srv := futuresClient.NewCreateOrderService()
	srv.Symbol(symbol)
	srv.Side(futures.SideTypeSell)
	srv.PositionSide(futures.PositionSideTypeLong)
	srv.Type(futures.OrderTypeStopMarket)
	srv.StopPrice("0.2")
	srv.ClosePosition(true)
	res, err := srv.Do(ctx)
	if err != nil {
		panic(err)
	}
	printVal(res)
}

func CloseTradeTrail(ctx context.Context, offer *Offer, tradeType TradeType) {
	if !apiSwitch {
		return
	}
	srv := futuresClient.NewCreateOrderService()
	srv.Symbol(offer.Symbol)
	switch tradeType {
	case TradeBuy:
		srv.Side(futures.SideTypeBuy)
		srv.PositionSide(futures.PositionSideTypeShort)
	case TradeSell:
		srv.Side(futures.SideTypeSell)
		srv.PositionSide(futures.PositionSideTypeLong)
	}
	srv.Type(futures.OrderTypeTrailingStopMarket)
	srv.Type(futures.OrderTypeTrailingStopMarket)
	srv.ActivationPrice(offer.Price)
	d := decimal.NewFromFloat(float64(offer.CallbackRate) / float64(tradeMargin))
	if d.LessThan(decimal.NewFromFloat(0.2)) {
		d = decimal.NewFromFloat(0.2)
	}
	srv.CallbackRate(d.Round(1).String())
	srv.Quantity(offer.Amount)
	res, err := srv.Do(ctx)
	if err != nil {
		panic(err)
	}
	printVal(res)
}

// case buy, with sell type, with long or short
func CloseTrade(ctx context.Context, offer *Offer, tradeType TradeType, closeType CloseType) {
	if !apiSwitch {
		return
	}
	srv := futuresClient.NewCreateOrderService()
	srv.Symbol(offer.Symbol)
	switch tradeType {
	case TradeBuy:
		srv.Side(futures.SideTypeBuy)
		srv.PositionSide(futures.PositionSideTypeShort)
	case TradeSell:
		srv.Side(futures.SideTypeSell)
		srv.PositionSide(futures.PositionSideTypeLong)
	}
	switch closeType {
	case CloseProfit:
		srv.Type(futures.OrderTypeTakeProfitMarket)
	case CloseStop:
		srv.Type(futures.OrderTypeStopMarket)
	}
	srv.StopPrice(offer.Price).ClosePosition(true)
	res, err := srv.Do(ctx)
	if err != nil {
		panic(err)
	}
	printVal(res)
}

func printVal(v interface{}) {
	b, _ := json.Marshal(v)
	fmt.Println(string(b))
}
