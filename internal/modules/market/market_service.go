package market

import (
	"context"

)

type MarketService struct {
	repo *MarketRepository
	hub *Hub
}

func NewMarketService(repo *MarketRepository, hub *Hub) *MarketService {
	return &MarketService{repo: repo, hub: hub}
}

func (s *MarketService) Publish(ctx context.Context, ticker Ticker) error {

	if err := s.repo.SaveTicker(ctx, ticker); err != nil {
		return err
	}

	s.hub.Broadcast(ticker.Symbol, ticker)

	return nil
}