package market

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type MarketRepository struct {
	rdb *redis.Client
}

func NewMarketRepo(rdb *redis.Client) *MarketRepository {
	return &MarketRepository{rdb: rdb}
}

func (r *MarketRepository) SaveTicker(ctx context.Context, ticker Ticker) error {

	key := "market:ticker:" + ticker.Symbol

	b, err := json.Marshal(ticker)
	if err != nil {
		return err
	}

	return r.rdb.Set(ctx, key, b, 0).Err()
} 

func (r *MarketRepository) GetTicker(ctx context.Context, symbol string) (*Ticker, error) {

	key := "market:ticker:" + symbol

	val, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var ticker Ticker
	err = json.Unmarshal([]byte(val), &ticker)

	return &ticker, err
}