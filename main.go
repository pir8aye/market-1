package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/lib/pq"
	"github.com/urfave/cli"
)

const (
	coinMarketCapURL = "https://api.coinmarketcap.com/v1/ticker/?convert=CNY&timestamp=%s"
)

var (
	min   = int64(60)
	min5  = 5 * min
	min10 = 2 * min5
	min15 = 3 * min5
	min30 = 2 * min15
	hour  = 2 * min30
	day   = 24 * hour
	week  = 7 * day
)

const (
	coinmarketcapmin     = "coinmarketcapmin"
	coinmarketcap5min    = "coinmarketcap5min"
	coinmarketcap10min   = "coinmarketcap10min"
	coinmarketcap30min   = "coinmarketcap30min"
	coinmarketcap15min   = "coinmarketcap15min"
	coinmarketcapday     = "coinmarketcapday"
	coinmarketcaphour    = "coinmarketcaphour"
	coinmarketcapweek    = "coinmarketcapweek"
	coinmarketcapcurrent = "coinmarketcapcurrent"
)

var table = "coinmarketcap"

// CREATE INDEX IF NOT EXISTS  pg 9.5 才支持
const tblCoinMarketCapXmin = `
	CREATE TABLE IF NOT EXISTS %s (
		id SERIAL PRIMARY KEY,
		asset_id character varying(32) NOT NULL,  
		name character varying(32) NOT NULL, 
		symbol character varying(32) NOT NULL, 
		rank integer  NOT NULL, 
		price_usd_first real  NOT NULL, 
		price_usd_last real NOT NULL, 
		price_usd_low real  NOT NULL, 
		price_usd_high  real NOT NULL, 
		price_btc_first real  NOT NULL, 
		price_btc_last real  NOT NULL, 
		price_btc_low real  NOT NULL, 
		price_btc_high real  NOT NULL, 
		price_cny_first real  NOT NULL, 
		price_cny_last real  NOT NULL, 
		price_cny_low real  NOT NULL, 
		price_cny_high real  NOT NULL, 
		last_updated bigint NOT NULL,
		timestamp bigint NOT NULL,
        _group character varying(32) 
	); 

   CREATE INDEX IF NOT EXISTS index_timestamp_%s ON %s (timestamp);
   CREATE INDEX IF NOT EXISTS index_last_updated_%s ON %s (last_updated);
   CREATE INDEX IF NOT EXISTS index_symbol_%s ON %s USING hash (symbol);
`

const tblCoinMarketCap = `
	CREATE TABLE IF NOT EXISTS %s (
		id SERIAL PRIMARY KEY,
		asset_id  character varying(32) NOT NULL,  
		name character varying(32) NOT NULL,
		symbol character varying(32) NOT NULL,
		rank character varying(32) NOT NULL,
		price_usd character varying(64) NOT NULL,
		price_btc character varying(64) NOT NULL,
		volume_usd_24h character varying(64) NOT NULL,
		market_cap_usd character varying(64) NOT NULL,
		available_supply character varying(64) NOT NULL,
		total_supply character varying(64) NOT NULL,
		percent_change_1h character varying(64) NOT NULL,
		percent_change_24h character varying(64) NOT NULL,
		percent_change_7d character varying(64) NOT NULL,
		last_updated character varying(64) NOT NULL,
		price_cny character varying(64) NOT NULL,
		volume_cny_24h character varying(64) NOT NULL,
		market_cap_cny character varying(64) NOT NULL 
	);
`

type priceCoinMarketCap struct {
	Id                string `json:"id"`
	Name              string `json:"name"`
	Symbol            string `json:"symbol"`
	Rank              string `json:"rank"`
	PriceUSD          string `json:"price_usd"`
	PriceBTC          string `json:"price_btc"`
	VolumeUSD24H      string `json:"24h_volume_usd"`
	MarketCapUSD      string `json:"market_cap_usd"`
	AvailableSupply   string `json:"available_supply"`
	TotalSupply       string `json:"total_supply"`
	PercentChanage1H  string `json:"percent_change_1h"`
	PercentChanage24H string `json:"percent_change_24h"`
	PercentChanage7D  string `json:"percent_change_7d"`
	LastUpdated       string `json:"last_updated"`
	PriceCNY          string `json:"price_cny"`
	VolumeCNY24H      string `json:"24h_volume_cny"`
	MarketCapCNY      string `json:"market_cap_cny"`
}

type kPriceCoinMarketCap struct {
	Id            string    `json:"id"`
	Name          string    `json:"name"`
	Symbol        string    `json:"symbol"`
	Rank          big.Int   `json:"rank"`
	PriceUSDFirst big.Float `json:"price_usd_first"`
	PriceUSDLast  big.Float `json:"price_usd_last"`
	PriceUSDLow   big.Float `json:"price_usd_low"`
	PriceUSDHigh  big.Float `json:"price_usd_high"`
	PriceBTCFirst big.Float `json:"price_btc_first"`
	PriceBTCLast  big.Float `json:"price_btc_last"`
	PriceBTCLow   big.Float `json:"price_btc_low"`
	PriceBTCHigh  big.Float `json:"price_btc_high"`
	PriceCNYFirst big.Float `json:"price_cny_first"`
	PriceCNYLast  big.Float `json:"price_cny_last"`
	PriceCNYLow   big.Float `json:"price_cny_low"`
	PriceCNYHigh  big.Float `json:"price_cny_high"`
	LastUpdated   big.Int   `json:"last_updated"`
}

type kPriceCoinMarketCapList struct {
	list      []kPriceCoinMarketCap
	timestamp int64
}

func (this *kPriceCoinMarketCapList) Copy() *kPriceCoinMarketCapList {
	ret := new(kPriceCoinMarketCapList)
	ret.timestamp = this.timestamp
	ret.list = append(ret.list, this.list...)
	return ret
}

func newKPriceCoinMarketCapList(list []priceCoinMarketCap, now int64) *kPriceCoinMarketCapList {
	ret := new(kPriceCoinMarketCapList)
	ret.timestamp = now
	for _, v := range list {
		var tmp kPriceCoinMarketCap
		tmp.Id = v.Id
		tmp.Name = v.Name
		tmp.Symbol = v.Symbol
		tmp.Rank.SetString(v.Rank, 10)
		tmp.PriceUSDFirst.SetString(v.PriceUSD)
		tmp.PriceBTCFirst.SetString(v.PriceBTC)
		tmp.PriceCNYFirst.SetString(v.PriceCNY)
		tmp.PriceUSDLast.SetString(v.PriceUSD)
		tmp.PriceBTCLast.SetString(v.PriceBTC)
		tmp.PriceCNYLast.SetString(v.PriceCNY)
		tmp.PriceUSDLow.SetString(v.PriceUSD)
		tmp.PriceBTCLow.SetString(v.PriceBTC)
		tmp.PriceCNYLow.SetString(v.PriceCNY)
		tmp.PriceUSDHigh.SetString(v.PriceUSD)
		tmp.PriceBTCHigh.SetString(v.PriceBTC)
		tmp.PriceCNYHigh.SetString(v.PriceCNY)
		tmp.LastUpdated.SetString(v.LastUpdated, 10)
		ret.list = append(ret.list, tmp)
	}
	return ret
}

func main() {
	app := &cli.App{
		Name: "ethtx",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "drivername",
				Value: "postgres",
				Usage: "database type default postgres",
			},
			&cli.StringFlag{
				Name:  "dbport",
				Value: "5432",
				Usage: "db port ",
			},
			&cli.StringFlag{
				Name:  "dbname",
				Value: "postgres",
				Usage: "database name default postgres",
			},
			&cli.StringFlag{
				Name:  "dbusername",
				Value: "postgres",
				Usage: "database user name",
			},
			&cli.StringFlag{
				Name:  "dbpassword",
				Value: "postgres",
				Usage: "database password of username",
			},
			&cli.StringFlag{
				Name:  "dbhost",
				Value: "localhost",
				Usage: "database host",
			},
			&cli.StringFlag{
				Name:  "tblname",
				Value: "pricecoinmarketcap",
				Usage: "market table",
			},
			&cli.DurationFlag{
				Name:  "interval",
				Value: 10 * time.Second,
				Usage: "price query duration",
			},
		},
		Action: func(c *cli.Context) error {
			// db config
			drivername := c.String("drivername")
			dbusername := c.String("dbusername")
			dbname := c.String("dbname")
			dbpassword := c.String("dbpassword")
			dbhost := c.String("dbhost")
			dbport := c.String("dbport")
			interval := c.Duration("interval")
			tblname := c.String("tblname")

			log.Println("drivername:", drivername)
			log.Println("dbusername:", dbusername)
			log.Println("dbhost:", dbhost)
			log.Println("dbname:", dbname)
			log.Println("dbpassword:", dbpassword)
			log.Println("dbport:", dbport)
			log.Println("dbusername", dbusername)
			log.Println("interval", interval)
			log.Println("tblname", tblname)
			table = tblname

			db, err := sql.Open(drivername, fmt.Sprintf("user=%v password=%v host=%v dbname=%v port=%v sslmode=disable", dbusername, dbpassword, dbhost, dbname, dbport))
			checkErr(err)
			ch := make(chan *kPriceCoinMarketCapList, 100)
			rt := make(chan *kPriceCoinMarketCapList, 100)
			timer := time.NewTicker(interval)
			checkErr(createTable(db, fmt.Sprintf(tblCoinMarketCap, tblname)))
			go dispatch(db, ch, 100)
			go realTimeAggregation(db, tblname, time.Now().Unix(), int64(interval), rt)
			for t := range timer.C {
				if resp, err := http.Get(fmt.Sprintf(coinMarketCapURL, t.Unix())); err == nil {
					var list []priceCoinMarketCap
					dec := json.NewDecoder(resp.Body)
					if err := dec.Decode(&list); err != nil {
						log.Println("decode error ", err)
						continue
					}
					x := newKPriceCoinMarketCapList(list, t.Unix())
					ch <- x.Copy()
					rt <- x.Copy()
					// 更新实时数据
					go func(x *kPriceCoinMarketCapList, list []priceCoinMarketCap) {
						if err := insert(db, tblname, list); err != nil {
							log.Println(err)
						}
					}(x, list)

				} else {
					log.Println(err)
				}
			}
			return nil
		},
	}
	app.Run(os.Args)
}

func dispatch(db *sql.DB, ch <-chan *kPriceCoinMarketCapList, size int) {
	var ch1 = make(chan *kPriceCoinMarketCapList, size)
	var ch2 = make(chan *kPriceCoinMarketCapList, size)
	var ch3 = make(chan *kPriceCoinMarketCapList, size)
	var ch4 = make(chan *kPriceCoinMarketCapList, size)
	var ch5 = make(chan *kPriceCoinMarketCapList, size)
	var ch6 = make(chan *kPriceCoinMarketCapList, size)
	var ch7 = make(chan *kPriceCoinMarketCapList, size)

	/*
		var r1 = make(chan *kPriceCoinMarketCapList, size)
		var r2 = make(chan *kPriceCoinMarketCapList, size)
		var r3 = make(chan *kPriceCoinMarketCapList, size)
		var r4 = make(chan *kPriceCoinMarketCapList, size)
		var r5 = make(chan *kPriceCoinMarketCapList, size)
		var r6 = make(chan *kPriceCoinMarketCapList, size)
		var r7 = make(chan *kPriceCoinMarketCapList, size)
		var r8 = make(chan *kPriceCoinMarketCapList, size)
	*/

	tm := time.Now()
	now := tm.Unix()
	checkErr(createTable(db, fmt.Sprintf(tblCoinMarketCapXmin, coinmarketcapmin,
		coinmarketcapmin, coinmarketcapmin,
		coinmarketcapmin, coinmarketcapmin,
		coinmarketcapmin, coinmarketcapmin)))
	checkErr(createTable(db, fmt.Sprintf(tblCoinMarketCapXmin, coinmarketcap5min,
		coinmarketcap5min, coinmarketcap5min,
		coinmarketcap5min, coinmarketcap5min,
		coinmarketcap5min, coinmarketcap5min)))
	checkErr(createTable(db, fmt.Sprintf(tblCoinMarketCapXmin, coinmarketcap10min,
		coinmarketcap10min, coinmarketcap10min,
		coinmarketcap10min, coinmarketcap10min,
		coinmarketcap10min, coinmarketcap10min)))
	checkErr(createTable(db, fmt.Sprintf(tblCoinMarketCapXmin, coinmarketcap30min,
		coinmarketcap30min, coinmarketcap30min,
		coinmarketcap30min, coinmarketcap30min,
		coinmarketcap30min, coinmarketcap30min)))
	checkErr(createTable(db, fmt.Sprintf(tblCoinMarketCapXmin, coinmarketcap15min,
		coinmarketcap15min, coinmarketcap15min,
		coinmarketcap15min, coinmarketcap15min,
		coinmarketcap15min, coinmarketcap15min)))
	checkErr(createTable(db, fmt.Sprintf(tblCoinMarketCapXmin, coinmarketcapday,
		coinmarketcapday, coinmarketcapday,
		coinmarketcapday, coinmarketcapday,
		coinmarketcapday, coinmarketcapday)))
	checkErr(createTable(db, fmt.Sprintf(tblCoinMarketCapXmin, coinmarketcaphour,
		coinmarketcaphour, coinmarketcaphour,
		coinmarketcaphour, coinmarketcaphour,
		coinmarketcaphour, coinmarketcaphour)))
	checkErr(createTable(db, fmt.Sprintf(tblCoinMarketCapXmin, coinmarketcapweek,
		coinmarketcapweek, coinmarketcapweek,
		coinmarketcapweek, coinmarketcapweek,
		coinmarketcapweek, coinmarketcapweek)))
	checkErr(createTable(db, fmt.Sprintf(tblCoinMarketCapXmin, coinmarketcapcurrent,
		coinmarketcapcurrent, coinmarketcapcurrent,
		coinmarketcapcurrent, coinmarketcapcurrent,
		coinmarketcapcurrent, coinmarketcapcurrent)))

	//  数据全部由一分钟数据出减少等待误差
	go aggregation(db, now-now%min+min, min, coinmarketcapmin, ch, ch1, ch2, ch3, ch4, ch5, ch6, ch7) // r1, r2, r3, r4, r5, r6, r7, r8)

	// go realTimeAggregation(db, coinmarketcapmin, now-now%min+min, min, r1)

	go aggregation(db, now+min5-now%min5, min5, coinmarketcap5min, ch1)
	// go realTimeAggregation(db, coinmarketcap5min, now-now%min5+min5, min5, r2)

	go aggregation(db, now+min10-now%min10, min10, coinmarketcap10min, ch2)
	// go realTimeAggregation(db, coinmarketcap10min, now-now%min10+min10, min10, r3)

	go aggregation(db, now+min15-now%min15, min15, coinmarketcap15min, ch3)
	// go realTimeAggregation(db, coinmarketcap15min, now-now%min15+min15, min15, r4)

	go aggregation(db, now+min30-now%min30, min30, coinmarketcap30min, ch4)
	// go realTimeAggregation(db, coinmarketcap30min, now-now%min30+min30, min30, r5)

	go aggregation(db, now+hour-now%hour, hour, coinmarketcaphour, ch5)
	// go realTimeAggregation(db, coinmarketcaphour, now-now%hour+hour, hour, r6)

	go aggregation(db, now+day-now%day, day, coinmarketcapday, ch6)
	// go realTimeAggregation(db, coinmarketcapday, now-now%day+day, day, r7)

	weekday := tm.Weekday()
	if weekday == time.Sunday {
		// 星期天统计星期天这天的就行了
		go aggregation(db, now+(day-now%day), week, coinmarketcapweek, ch7)
		//	go realTimeAggregation(db, coinmarketcapweek, now-now%day+day, week, r7)
	} else {
		x := int64(time.Saturday-weekday) + 2
		go aggregation(db, now+day-(now%day)+x*day, week, coinmarketcapweek, ch7)
		//	go realTimeAggregation(db, coinmarketcapweek, now-now%day+day+x*day, week, r7)
	}

	go func() {
		ticker := time.NewTicker(time.Hour)
		for {
			select {
			case now := <-ticker.C:
				exec(db, fmt.Sprintf("delete from %v where last_updated < $1;", table), fmt.Sprint(now.Unix()-7*24*3600))
			}
		}
	}()
}

func exec(db *sql.DB, sql string, args ...interface{}) error {
	stmt, err := db.Prepare(sql)
	if err != nil {
		log.Println("Prepare err:", err)
		return err
	}
	defer stmt.Close()
	if _, err = stmt.Exec(args...); err != nil {
		log.Println("del err:", err)
		return err
	}
	return nil
}

func summary(x, y *kPriceCoinMarketCapList) {
	for i := range x.list {
		for _, v := range y.list {
			if x.list[i].Symbol == v.Symbol {
				// 汇聚数据将
				x.list[i].PriceCNYLast = v.PriceCNYLast
				x.list[i].PriceBTCLast = v.PriceBTCLast
				x.list[i].PriceUSDLast = v.PriceUSDLast
				x.list[i].LastUpdated = v.LastUpdated
				x.list[i].Rank = v.Rank

				if x.list[i].PriceBTCLow.Cmp(&v.PriceBTCLow) > 0 {
					x.list[i].PriceBTCLow = v.PriceBTCLow
				}

				if x.list[i].PriceBTCHigh.Cmp(&v.PriceBTCHigh) < 0 {
					x.list[i].PriceBTCHigh = v.PriceBTCHigh
				}

				if x.list[i].PriceCNYLow.Cmp(&v.PriceCNYLow) > 0 {
					x.list[i].PriceCNYLow = v.PriceCNYLow
				}

				if x.list[i].PriceCNYHigh.Cmp(&v.PriceCNYHigh) < 0 {
					x.list[i].PriceCNYHigh = v.PriceCNYHigh
				}

				if x.list[i].PriceUSDLow.Cmp(&v.PriceUSDLow) > 0 {
					x.list[i].PriceUSDLow = v.PriceUSDLow
				}

				if x.list[i].PriceUSDHigh.Cmp(&v.PriceUSDHigh) < 0 {
					x.list[i].PriceUSDHigh = v.PriceUSDHigh
				}
			}
		}
	}
}

func upsertCoinMarketCapCurrent(db *sql.DB, tblname string, dat *kPriceCoinMarketCapList) error {
	return tx(db, func(txn *sql.Tx) error {
		if _, err := txn.Exec(fmt.Sprintf("delete from %v where _group = $1 ;", coinmarketcapcurrent), tblname); err != nil {
			return err
		}
		stmt, err := txn.Prepare(pq.CopyIn(coinmarketcapcurrent, "asset_id", "name", "symbol",
			"rank",
			"price_usd_first",
			"price_usd_last",
			"price_usd_low",
			"price_usd_high",
			"price_btc_first",
			"price_btc_last",
			"price_btc_low",
			"price_btc_high",
			"price_cny_first",
			"price_cny_last",
			"price_cny_low",
			"price_cny_high",
			"last_updated",
			"timestamp",
			"_group"))
		if err != nil {
			return err
		}

		for _, v := range dat.list {
			a1, _ := v.PriceUSDFirst.Float64()
			a2, _ := v.PriceUSDLast.Float64()
			a3, _ := v.PriceUSDLow.Float64()
			a4, _ := v.PriceUSDHigh.Float64()
			b1, _ := v.PriceBTCFirst.Float64()
			b2, _ := v.PriceBTCLast.Float64()
			b3, _ := v.PriceBTCLow.Float64()
			b4, _ := v.PriceBTCHigh.Float64()
			c1, _ := v.PriceCNYFirst.Float64()
			c2, _ := v.PriceCNYLast.Float64()
			c3, _ := v.PriceCNYLow.Float64()
			c4, _ := v.PriceCNYHigh.Float64()
			_, err := stmt.Exec(v.Id,
				v.Name,
				v.Symbol,
				v.Rank.Int64(),
				a1, a2, a3, a4,
				b1, b2, b3, b4,
				c1, c2, c3, c4,
				v.LastUpdated.Int64(),
				dat.timestamp,
				tblname)
			if err != nil {
				return err
			}
		}
		_, err = stmt.Exec()
		if err != nil {
			return err
		}
		err = stmt.Close()
		if err != nil {
			return err
		}
		return nil
	})

}

// 实时数据汇总
func realTimeAggregation(db *sql.DB, tblname string, base, interval int64, in <-chan *kPriceCoinMarketCapList) {
	var current *kPriceCoinMarketCapList
	var tmp *kPriceCoinMarketCapList
	var resetCurrent = base
	var resetCurrentInterval = interval
	// 实时数据重置 至少单位为一天
	if interval < day {
		resetCurrent = base - base%day + day
		resetCurrentInterval = day
	}

	ticker := time.NewTicker(time.Second)
	for {
		select {
		case t := <-ticker.C:
			// 更新first
			if tmp == nil {
				continue
			}

			if current == nil {
				continue
			}

			if t.Unix() >= base {
				base += interval
				summary(tmp, current)
				current = tmp
			}

			if t.Unix() < resetCurrent {
				continue
			}
			resetCurrent += resetCurrentInterval
			current = nil
		case x := <-in: // 一旦有数据变动更新
			tmp = x.Copy()
			if current == nil {
				current = x
			} else {
				summary(current, x)
			}
			if err := upsertCoinMarketCapCurrent(db, tblname, current); err != nil {
				log.Println("tx error ", err)
			}
		}
	}
}

func aggregation(db *sql.DB, base, interval int64, tbl string, in <-chan *kPriceCoinMarketCapList, outs ...chan<- *kPriceCoinMarketCapList) {
	var tmp *kPriceCoinMarketCapList
	// 实时数据
	ticker := time.NewTicker(time.Second)
	next := base
	// 把误差控制在0s内
	for {
		select {
		case t := <-ticker.C: // beatheart
			if now := t.Unix(); now < next {
				continue
			}
			next += interval
			select {
			case v := <-in:
				if tmp == nil {
					tmp = v
				} else {
					summary(tmp, v) // 数据汇总
				}
			case <-time.After(3 * time.Second):
			}

			if tmp == nil {
				continue
			}
			if len(outs) > 0 {
				for _, ch := range outs {
					// dump
					ch <- tmp.Copy()
				}
			}

			now := time.Now()
			saveKPriceCoinMarketCap(db, tbl, "", tmp)
			fmt.Println(tbl, ":", time.Since(now))
			tmp = nil
		case v := <-in:
			if tmp == nil {
				tmp = v
				continue
			}
			summary(tmp, v)
		}
	}
}

func tx(db *sql.DB, fn func(tx *sql.Tx) error) error {
	txn, err := db.Begin()
	if err != nil {
		return err
	}
	if err := fn(txn); err != nil {
		txn.Rollback()
		return err
	}
	return txn.Commit()
}

func saveKPriceCoinMarketCap(db *sql.DB, tblname, group string, dat *kPriceCoinMarketCapList) error {
	txn, err := db.Begin()
	if err != nil {
		log.Println("save error", err)
		return err
	}
	stmt, err := txn.Prepare(pq.CopyIn(tblname, "asset_id", "name", "symbol",
		"rank",
		"price_usd_first",
		"price_usd_last",
		"price_usd_low",
		"price_usd_high",
		"price_btc_first",
		"price_btc_last",
		"price_btc_low",
		"price_btc_high",
		"price_cny_first",
		"price_cny_last",
		"price_cny_low",
		"price_cny_high",
		"last_updated",
		"timestamp",
		"_group"))
	if err != nil {
		txn.Rollback()
		log.Println("save error", err)
		return err
	}

	for _, v := range dat.list {
		a1, _ := v.PriceUSDFirst.Float64()
		a2, _ := v.PriceUSDLast.Float64()
		a3, _ := v.PriceUSDLow.Float64()
		a4, _ := v.PriceUSDHigh.Float64()
		b1, _ := v.PriceBTCFirst.Float64()
		b2, _ := v.PriceBTCLast.Float64()
		b3, _ := v.PriceBTCLow.Float64()
		b4, _ := v.PriceBTCHigh.Float64()
		c1, _ := v.PriceCNYFirst.Float64()
		c2, _ := v.PriceCNYLast.Float64()
		c3, _ := v.PriceCNYLow.Float64()
		c4, _ := v.PriceCNYHigh.Float64()
		_, err := stmt.Exec(v.Id,
			v.Name,
			v.Symbol,
			v.Rank.Int64(),
			a1, a2, a3, a4,
			b1, b2, b3, b4,
			c1, c2, c3, c4,
			v.LastUpdated.Int64(),
			dat.timestamp,
			group)
		if err != nil {
			txn.Rollback()
			log.Println(err, v)
			return err
		}
	}
	_, err = stmt.Exec()
	if err != nil {
		txn.Rollback()
		log.Println("save error", err)
	}

	err = stmt.Close()
	if err != nil {
		txn.Rollback()
		log.Fatal(err)
		return err
	}
	return txn.Commit()
}

func insert(db *sql.DB, tblname string, list []priceCoinMarketCap) error {
	txn, err := db.Begin()
	if err != nil {
		return err
	}
	// 事务
	stmt, err := txn.Prepare(pq.CopyIn(tblname,
		"asset_id",
		"name",
		"symbol",
		"rank",
		"price_usd",
		"price_btc",
		"volume_usd_24h",
		"market_cap_usd",
		"available_supply",
		"total_supply",
		"percent_change_1h",
		"percent_change_24h",
		"percent_change_7d",
		"last_updated",
		"price_cny",
		"volume_cny_24h",
		"market_cap_cny"))
	if err != nil {
		txn.Rollback()
		return err
	}

	for _, v := range list {
		_, err := stmt.Exec(v.Id, v.Name, v.Symbol, v.Rank, v.PriceUSD, v.PriceBTC, v.VolumeUSD24H, v.MarketCapUSD, v.AvailableSupply,
			v.TotalSupply, v.PercentChanage1H, v.PercentChanage24H, v.PercentChanage7D, v.LastUpdated, v.PriceCNY, v.VolumeCNY24H, v.MarketCapCNY)

		if err != nil {
			txn.Rollback()
			log.Println(err, v)
			return err
		}
	}
	_, err = stmt.Exec()
	if err != nil {
		txn.Rollback()
		log.Fatal(err)
	}

	err = stmt.Close()
	if err != nil {
		txn.Rollback()
		return err
	}
	return txn.Commit()
}

func createTable(db *sql.DB, sql string) error {
	_, err := db.Exec(sql)
	_ = err
	return nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
