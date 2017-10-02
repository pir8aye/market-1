package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"sort"

	"github.com/urfave/cli"
)

var coinMarketCapURL = "https://api.coinmarketcap.com/v1/ticker/?convert=CNY"

type coinMarketCap struct {
	Symbol       string     `json:"symbol"`
	Rank         string     `json:"rank"`
	PriceUSD     *big.Float `json:"price_usd"`
	MarketCapUSD *big.Float `json:"market_cap_usd"`
}

func coinMarketCapTopN(n int) []coinMarketCap {
	resp, err := http.Get(coinMarketCapURL)
	if err != nil {
		log.Fatal("err", err)
	}
	var list []coinMarketCap
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&list)
	// 数据量太小没必要用堆
	sort.Slice(list, func(i, j int) bool {
		if list[i].MarketCapUSD == nil {
			return false
		}

		if list[j].MarketCapUSD == nil {
			return true
		}
		return list[i].MarketCapUSD.Cmp(list[j].MarketCapUSD) > 0
	})
	if len(list) > n {
		list = list[:n]
	}
	return list
}

func main() {
	app := &cli.App{
		Name: "print",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "size",
				Value: 100,
				Usage: "size of top",
			}},
		Action: func(c *cli.Context) error {
			size := c.Int("size")
			topN := coinMarketCapTopN(size)
			for i, v := range topN {
				fmt.Printf("%v,%v,%f,%f \n", i+1, v.Symbol, v.MarketCapUSD, v.PriceUSD)
			}
			return nil
		},
	}
	app.Run(os.Args)
}
