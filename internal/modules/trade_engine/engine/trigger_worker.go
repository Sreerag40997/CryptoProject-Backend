package engine

import (
	"context"
	"cryptox/internal/modules/trade_engine/model"
	"encoding/json"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type TriggerRepository interface {
	GetOpenTriggerOrders(ctx context.Context) ([]model.Order, error)
	UpdateOrder(ctx context.Context, order *model.Order) error
}

type Ticker struct {
	LastPrice string `json:"lastPrice"`
}

type TriggerWatcher struct {
	engine *Engine
	repo   TriggerRepository
	redis  *redis.Client
}

func NewTriggerWatcher(
	engine *Engine,
	repo TriggerRepository,
	redis *redis.Client,
) *TriggerWatcher {

	return &TriggerWatcher{
		engine: engine,
		repo:   repo,
		redis:  redis,
	}
}

func (tw *TriggerWatcher) Start() {

	go func() {

		for {

			orders, err := tw.repo.GetOpenTriggerOrders(
				context.Background(),
			)

			if err != nil {
				time.Sleep(time.Second)
				continue
			}

			for _, order := range orders {

				currentPrice := tw.getMarketPrice(order.Symbol)

				if currentPrice == 0 {
					continue
				}

				// STOP LOSS

				if order.Type == "stop_loss" {

					if currentPrice <= order.StopPrice {

						order.Type = "market"
						order.Price = currentPrice

						tw.engine.Submit(&order)
					}
				}

				// TAKE PROFIT

				if order.Type == "take_profit" {

					if currentPrice >= order.TargetPrice {

						order.Type = "market"
						order.Price = currentPrice

						tw.engine.Submit(&order)
					}
				}
			}

			time.Sleep(time.Second)
		}
	}()
}

func (tw *TriggerWatcher) getMarketPrice(symbol string) int64 {

	key := "market:ticker:" + symbol

	val, err := tw.redis.Get(
		context.Background(),
		key,
	).Result()

	if err != nil {
		return 0
	}

	var ticker Ticker

	err = json.Unmarshal([]byte(val), &ticker)
	if err != nil {
		return 0
	}

	price, _ := strconv.ParseInt(
		ticker.LastPrice,
		10,
		64,
	)

	return price
}