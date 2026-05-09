package bot

import (
	"context"
	"cryptox/internal/modules/trade_engine/engine"
	"cryptox/internal/modules/trade_engine/model"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Bot struct {
	Engine *engine.Engine
	Redis  *redis.Client
}

func (b *Bot) Start() {

	go func() {
		for {

			price := b.getMarketPrice("BTC-INR")
			if price == 0 {
				time.Sleep(time.Second)
				continue
			}

			b.Engine.Submit(&model.Order{
				UserID:   999,
				Symbol:   "BTC-INR",
				Side:     "buy",
				Type:     "limit",
				Price:    price - 50,
				Quantity: 1,
			})

			b.Engine.Submit(&model.Order{
				UserID:   999,
				Symbol:   "BTC-INR",
				Side:     "sell",
				Type:     "limit",
				Price:    price + 50,
				Quantity: 1,
			})

			time.Sleep(2 * time.Second)
		}
	}()
}

func (b *Bot) getMarketPrice(symbol string) int64 {

	key := "market:price:" + symbol

	val, err := b.Redis.Get(context.Background(), key).Result()
	if err != nil {
		return 0
	}

	price, _ := strconv.ParseInt(val, 10, 64)
	return price
}