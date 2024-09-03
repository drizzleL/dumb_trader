package binanceapi

import (
	"context"
	"encoding/json"
	"log"

	"github.com/adshao/go-binance/v2/futures"
)

func GetPosition(ctx context.Context) {
	acc, err := futuresClient.NewGetAccountService().Do(ctx)
	if err != nil {
		panic(err)
	}
	for _, p := range acc.Positions {
		if p.Symbol != "WLDUSDT" {
			continue
		}
		if p.PositionAmt == "0" {
			continue
		}
		if p.PositionSide == futures.PositionSideTypeLong {

		}
		b, _ := json.Marshal(p)
		log.Println(string(b))
	}
}
