package engine

import (
	"cryptox/internal/modules/trade_engine/model"
	"sort"
	"sync"
)

type OrderBook struct {
	Bids map[string]map[int64][]*model.Order
	Asks map[string]map[int64][]*model.Order

	BidPrices map[string][]int64
	AskPrices map[string][]int64

	mu sync.RWMutex
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		Bids:      make(map[string]map[int64][]*model.Order),
		Asks:      make(map[string]map[int64][]*model.Order),
		BidPrices: make(map[string][]int64),
		AskPrices: make(map[string][]int64),
	}
}

func (ob *OrderBook) AddOrder(order *model.Order) {

	ob.mu.Lock()
	defer ob.mu.Unlock()

	symbol := order.Symbol

	if order.Side == "buy" {

		if ob.Bids[symbol] == nil {
			ob.Bids[symbol] = make(map[int64][]*model.Order)
		}

		if _, ok := ob.Bids[symbol][order.Price]; !ok {
			ob.BidPrices[symbol] = append(ob.BidPrices[symbol], order.Price)

			sort.Slice(ob.BidPrices[symbol], func(i, j int) bool {
				return ob.BidPrices[symbol][i] > ob.BidPrices[symbol][j]
			})
		}

		ob.Bids[symbol][order.Price] =
			append(ob.Bids[symbol][order.Price], order)

		return
	}

	// SELL SIDE

	if ob.Asks[symbol] == nil {
		ob.Asks[symbol] = make(map[int64][]*model.Order)
	}

	if _, ok := ob.Asks[symbol][order.Price]; !ok {
		ob.AskPrices[symbol] = append(ob.AskPrices[symbol], order.Price)

		sort.Slice(ob.AskPrices[symbol], func(i, j int) bool {
			return ob.AskPrices[symbol][i] < ob.AskPrices[symbol][j]
		})
	}

	ob.Asks[symbol][order.Price] =
		append(ob.Asks[symbol][order.Price], order)
}

func (ob *OrderBook) Remove(order *model.Order) {

	ob.mu.Lock()
	defer ob.mu.Unlock()

	symbol := order.Symbol

	if order.Side == "buy" {

		orders := ob.Bids[symbol][order.Price]

		var updated []*model.Order

		for _, o := range orders {
			if o.ID != order.ID {
				updated = append(updated, o)
			}
		}

		ob.Bids[symbol][order.Price] = updated

		return
	}

	orders := ob.Asks[symbol][order.Price]

	var updated []*model.Order

	for _, o := range orders {
		if o.ID != order.ID {
			updated = append(updated, o)
		}
	}

	ob.Asks[symbol][order.Price] = updated
}