package engine

import (
	"context"
	"cryptox/internal/modules/trade_engine/model"
)

type Engine struct {
	orderBook *OrderBook
	executor  *Executor
	queue     chan *model.Order
}

func NewEngine(executor *Executor) *Engine {
	return &Engine{
		orderBook: NewOrderBook(),
		executor:  executor,
		queue:     make(chan *model.Order, 1000),
	}
}

func (e *Engine) Start() {
	go func() {
		for order := range e.queue {
			e.process(order)
		}
	}()
}

func (e *Engine) Submit(order *model.Order) {
	e.queue <- order
}

func (e *Engine) process(order *model.Order) {

	matches := e.orderBook.Match(order)

	for _, match := range matches {

		if order.RemainingQty == 0 {
			break
		}

		qty := min(order.RemainingQty, match.RemainingQty)

		err := e.executor.Execute(
			context.Background(),
			order,
			match,
			qty,
			match.Price,
		)

		if err != nil {
			continue
		}
	}

	if order.RemainingQty > 0 {
		e.orderBook.AddOrder(order)
	}
}

