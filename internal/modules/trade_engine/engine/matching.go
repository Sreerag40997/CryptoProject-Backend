package engine

import "cryptox/internal/modules/trade_engine/model"

func (ob *OrderBook) Match(order *model.Order) []*model.Order {

	ob.mu.Lock()
	defer ob.mu.Unlock()

	var result []*model.Order

	if order.Side == "buy" {

		for _, price := range ob.AskPrices {

			if price > order.Price {
				break
			}

			orders := ob.Asks[price]

			for _, o := range orders {

				if order.RemainingQty == 0 {
					break
				}

				result = append(result, o)
			}
		}
	}

	// mirror for sell

	return result
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
