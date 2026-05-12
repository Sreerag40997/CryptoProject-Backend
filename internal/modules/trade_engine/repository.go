package tradeengine

import (
	"context"
	"cryptox/internal/modules/trade_engine/model"

	"gorm.io/gorm"
)


type Repository interface {
	// user
	CreateOrder(ctx context.Context, order *model.Order) error
	GetOrderByID(ctx context.Context, id uint) (*model.Order, error)
	GetOrdersByUser(ctx context.Context, userID uint, status string, limit, offset int) ([]model.Order, error)

	CancelOrder(ctx context.Context, order *model.Order) error

	GetOpenOrdersForMatching(ctx context.Context, symbol, side string) ([]model.Order, error)

	// execution
	UpdateOrder(ctx context.Context, order *model.Order) error
	CreateTrade(ctx context.Context, trade *model.Trade) error
	CreateOrderFill(ctx context.Context, fill *model.OrderFill) error

	// admin
	GetAllOrders(ctx context.Context, limit, offset int) ([]model.Order, error)
	GetAllTrades(ctx context.Context, limit, offset int) ([]model.Trade, error)

	// transaction
	WithTx(ctx context.Context, fn func(Repository) error) error

  GetTradesByUser(ctx context.Context, userID uint, limit, offset int) ([]model.Trade, error)
  GetOrderFills(ctx context.Context, orderID uint) ([]model.OrderFill, error)
  GetTradesBySymbol(ctx context.Context, symbol string, limit int) ([]model.Trade, error)
  GetOpenTriggerOrders(ctx context.Context) ([]model.Order, error)
}


type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{db: db}
}

// user side

func (r *repo) CreateOrder(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *repo) GetOrderByID(ctx context.Context, id uint) (*model.Order, error) {
	var order model.Order
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&order).Error

	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *repo) GetOrdersByUser(ctx context.Context, userID uint, status string, limit, offset int) ([]model.Order, error) {
	var orders []model.Order

	query := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("id DESC").
		Limit(limit).
		Offset(offset)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Find(&orders).Error
	return orders, err
}

func (r *repo) CancelOrder(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).
		Model(order).
		Update("status", "cancelled").Error
}
func (r *repo) GetTradesByUser(ctx context.Context, userID uint, limit, offset int) ([]model.Trade, error) {
	var trades []model.Trade

	err := r.db.WithContext(ctx).
		Where("buyer_id = ? OR seller_id = ?", userID, userID).
		Order("id DESC").
		Limit(limit).
		Offset(offset).
		Find(&trades).Error

	return trades, err
}

// matching engine

func (r *repo) GetOpenOrdersForMatching(ctx context.Context, symbol, side string) ([]model.Order, error) {
	var orders []model.Order

	query := r.db.WithContext(ctx).
		Where("symbol = ? AND side = ? AND status = ?", symbol, side, "open")

	// buy-- match sell (lowest first)
	if side == "sell" {
		query = query.Order("price ASC, created_at ASC")
	}

	// sell -- match buy (highest first)
	if side == "buy" {
		query = query.Order("price DESC, created_at ASC")
	}

	err := query.Find(&orders).Error
	return orders, err
}

//execution

func (r *repo) UpdateOrder(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).
		Save(order).Error
}

func (r *repo) CreateTrade(ctx context.Context, trade *model.Trade) error {
	return r.db.WithContext(ctx).
		Create(trade).Error
}

func (r *repo) CreateOrderFill(ctx context.Context, fill *model.OrderFill) error {
	return r.db.WithContext(ctx).
		Create(fill).Error
}

func (r *repo) GetOrderFills(ctx context.Context, orderID uint) ([]model.OrderFill, error) {
	var fills []model.OrderFill

	err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Order("id ASC").
		Find(&fills).Error

	return fills, err
}


func (r *repo) GetTradesBySymbol(ctx context.Context, symbol string, limit int) ([]model.Trade, error) {
	var trades []model.Trade

	err := r.db.WithContext(ctx).
		Where("symbol = ?", symbol).
		Order("id DESC").
		Limit(limit).
		Find(&trades).Error

	return trades, err
}

//admin side

func (r *repo) GetAllOrders(ctx context.Context, limit, offset int) ([]model.Order, error) {
	var orders []model.Order

	err := r.db.WithContext(ctx).
		Order("id DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error

	return orders, err
}

func (r *repo) GetAllTrades(ctx context.Context, limit, offset int) ([]model.Trade, error) {
	var trades []model.Trade

	err := r.db.WithContext(ctx).
		Order("id DESC").
		Limit(limit).
		Offset(offset).
		Find(&trades).Error

	return trades, err
}

//transaction

func (r *repo) WithTx(ctx context.Context, fn func(Repository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &repo{db: tx}
		return fn(txRepo)
	})
}

func (r *repo) GetOpenTriggerOrders(
	ctx context.Context,
) ([]model.Order, error) {

	var orders []model.Order

	err := r.db.WithContext(ctx).
		Where(
			"(type = ? OR type = ?) AND status = ?",
			"stop_loss",
			"take_profit",
			"open",
		).
		Find(&orders).Error

	return orders, err
}
