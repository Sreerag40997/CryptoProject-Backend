package engine

import "cryptox/internal/modules/trade_engine/model"

func (ob *OrderBook) Match(order *model.Order) []*model.Order {

	ob.mu.RLock()
	defer ob.mu.RUnlock()

	var result []*model.Order

	symbol := order.Symbol

	// BUY ORDER
	if order.Side == "buy" {

		for _, price := range ob.AskPrices[symbol] {

			if price > order.Price {
				break
			}

			orders := ob.Asks[symbol][price]

			for _, o := range orders {

				if o.RemainingQty <= 0 {
					continue
				}

				result = append(result, o)
			}
		}

		return result
	}

	// SELL ORDER

	for _, price := range ob.BidPrices[symbol] {

		if price < order.Price {
			break
		}

		orders := ob.Bids[symbol][price]

		for _, o := range orders {

			if o.RemainingQty <= 0 {
				continue
			}

			result = append(result, o)
		}
	}

	return result
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}

	return b
}