package market

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)



func StartBinanceStream(market *MarketService) {

	url := "wss://stream.binance.com:9443/stream?streams=" +
		"btcusdt@ticker/" +
		"ethusdt@ticker/" +
		"solusdt@ticker/" +
		"bnbusdt@ticker/" +
		"xrpusdt@ticker/" +
		"adausdt@ticker/" +
		"dogeusdt@ticker/" +
		"dotusdt@ticker/" +
		"maticusdt@ticker/" +
		"linkusdt@ticker"

	for {

		// connect Binance websocket
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			log.Println("connect error:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		log.Println("binance stream connected")

		for {

			// read websocket message
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("stream disconnected:", err)
				conn.Close()
				break
			}

			// combined stream wrapper
			var wrap struct {
				Stream string        `json:"stream"`
				Data   BinanceTicker `json:"data"`
			}

			// decode safely
			decoder := json.NewDecoder(bytes.NewReader(msg))
			decoder.UseNumber()

			if err := decoder.Decode(&wrap); err != nil {
				log.Println("json error:", err)
				continue
			}

			// convert to your internal ticker DTO
			ticker := Ticker{
				Symbol:      wrap.Data.Symbol,
				LastPrice:   wrap.Data.LastPrice.String(),
				PriceChange: wrap.Data.PriceChange.String(),
				ChangePct:   wrap.Data.ChangePct.String(),
				Volume:      wrap.Data.Volume.String(),
				High:        wrap.Data.High.String(),
				Low:         wrap.Data.Low.String(),
				EventTime:   wrap.Data.EventTime.String(),
			}

			// publish to redis + websocket clients
			if err := market.Publish(context.Background(), ticker); err != nil {
				log.Println("publish error:", err)
			}
		}

		// reconnect delay
		time.Sleep(3 * time.Second)
	}
}