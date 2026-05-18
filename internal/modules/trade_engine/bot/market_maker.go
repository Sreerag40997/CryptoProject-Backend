package bot

import (
	"context"
	"cryptox/internal/modules/trade_engine/engine"
	"cryptox/internal/modules/trade_engine/model"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type Ticker struct {
	LastPrice string `json:"lastPrice"`
}

type Bot struct {
	Engine *engine.Engine
	Redis  *redis.Client
	Repo   engine.Repository
}

func (b *Bot) Start() {

	fmt.Println("market maker bot started")

	symbols := []string{
		"BTC-INR",
		"ETH-INR",
		"SOL-INR",
	}

	go func() {

		for {

			for _, symbol := range symbols {

				price := b.getMarketPrice(symbol)

				if price == 0 {
					continue
				}

				// BUY ORDER

				buyOrder := &model.Order{
					UserID:       999,
					Symbol:       symbol,
					Side:         "buy",
					Type:         "limit",
					Price:        price - 50,
					Quantity:     100000000,
					RemainingQty: 100000000,
					Status:       "open",
				}

				_ = b.Repo.CreateOrder(context.Background(), buyOrder)

				b.Engine.Submit(buyOrder)

				// SELL ORDER

				sellOrder := &model.Order{
					UserID:       999,
					Symbol:       symbol,
					Side:         "sell",
					Type:         "limit",
					Price:        price + 50,
					Quantity:     100000000,
					RemainingQty: 100000000,
					Status:       "open",
				}

				_ = b.Repo.CreateOrder(context.Background(), sellOrder)

				b.Engine.Submit(sellOrder)
			}

			time.Sleep(2 * time.Second)
		}
	}()
}

func (b *Bot) getMarketPrice(symbol string) int64 {

	cleanCoin := strings.Split(symbol, "-")[0]
	key := "market:ticker:" + cleanCoin + "USDT"

	val, err := b.Redis.Get(context.Background(), key).Result()
	if err != nil {
		return 0
	}

	var ticker Ticker

	err = json.Unmarshal([]byte(val), &ticker)
	if err != nil {
		return 0
	}

	rawPrice, err := strconv.ParseFloat(ticker.LastPrice, 64)
	if err != nil {
		return 0
	}

	return int64(rawPrice * 83.50 * 100)
}