package binanceapi

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/delivery"
	"github.com/adshao/go-binance/v2/futures"
)

var client *binance.Client
var futuresClient *futures.Client
var deliveryClient *delivery.Client

func Init(apiKey, secretKey string) {
	client = binance.NewClient(apiKey, secretKey)
	proxy := "socks5://127.0.0.1:1081"
	_proxy, _ := url.Parse(proxy)

	tr := &http.Transport{
		Proxy:           http.ProxyURL(_proxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	httpClient := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 5,
	}
	client.HTTPClient = httpClient

	futuresClient = binance.NewFuturesClient(apiKey, secretKey) // USDT-M Futures
	futuresClient.HTTPClient = httpClient

	deliveryClient = binance.NewDeliveryClient(apiKey, secretKey) // Coin-M Futures
	deliveryClient.HTTPClient = httpClient
}

func ShowPrice(ctx context.Context, symbol string) string {
	price := client.NewListPricesService()
	res, err := price.Symbol(symbol).Do(ctx)
	if err != nil {
		panic(err)
	}
	return res[0].Price
}
func ShowFuturesPrice(ctx context.Context, symbol string) []*futures.PriceChangeStats {
	srv := futuresClient.NewListPriceChangeStatsService()
	res, err := srv.Symbol(symbol).Do(ctx)
	if err != nil {
		panic(err)
	}
	return res
}
