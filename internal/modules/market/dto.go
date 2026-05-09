package market

import "encoding/json"

type Ticker struct {
	Symbol      string `json:"symbol"`
	LastPrice   string `json:"lastPrice"`
	PriceChange string `json:"priceChange"`
	ChangePct   string `json:"changePct"`
	Volume      string `json:"volume"`
	High        string `json:"high"`
	Low         string `json:"low"`
	EventTime   string `json:"eventTime"`
}

// Binance ticker payload
type BinanceTicker struct {
	EventType string      `json:"e"`
	EventTime json.Number `json:"E"`

	Symbol string `json:"s"`

	PriceChange json.Number `json:"p"`
	ChangePct   json.Number `json:"P"`

	WeightedAvg json.Number `json:"w"`

	PrevClose json.Number `json:"x"`

	LastPrice json.Number `json:"c"`

	LastQty json.Number `json:"Q"`

	BestBid    json.Number `json:"b"`
	BestBidQty json.Number `json:"B"`

	BestAsk    json.Number `json:"a"`
	BestAskQty json.Number `json:"A"`

	OpenPrice json.Number `json:"o"`

	High json.Number `json:"h"`
	Low  json.Number `json:"l"`

	Volume      json.Number `json:"v"`
	QuoteVolume json.Number `json:"q"`

	OpenTime  json.Number `json:"O"`
	CloseTime json.Number `json:"C"`

	FirstTradeID json.Number `json:"F"`
	LastTradeID  json.Number `json:"L"`

	TradeCount json.Number `json:"n"`
}