package engine

import (
	"cryptox/internal/modules/trade_engine/model"
	"sync"
)

type OrderBook struct {
	Bids map[int64][]*model.Order // price → orders
	Asks map[int64][]*model.Order

	BidPrices []int64
	AskPrices []int64

	mu sync.Mutex
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		Bids: make(map[int64][]*model.Order),
		Asks: make(map[int64][]*model.Order),
	}
}

func (ob *OrderBook) AddOrder(order *model.Order) {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	if order.Side == "buy" {
		ob.Bids[order.Price] = append(ob.Bids[order.Price], order)
	} else {
		ob.Asks[order.Price] = append(ob.Asks[order.Price], order)
	}
}