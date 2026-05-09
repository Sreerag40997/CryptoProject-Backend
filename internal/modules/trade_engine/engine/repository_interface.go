package engine

import (
	"context"
	"cryptox/internal/modules/trade_engine/model"
)

type Repository interface {
	CreateTrade(ctx context.Context, trade *model.Trade) error
	CreateOrderFill(ctx context.Context, fill *model.OrderFill) error
	UpdateOrder(ctx context.Context, order *model.Order) error

}

