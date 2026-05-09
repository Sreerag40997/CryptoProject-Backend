package tradeengine

import (
	"context"
	"cryptox/internal/modules/trade_engine/engine"
	"cryptox/internal/modules/trade_engine/model"
	"errors"
	"strconv"

	"github.com/redis/go-redis/v9"
)


type Service interface {

	// user side
	PlaceOrder(ctx context.Context, userID uint, req CreateOrderReq) error
	GetMyOrders(ctx context.Context, userID uint, status string, limit, offset int) ([]model.Order, error)
	GetOrderByID(ctx context.Context, userID uint, orderID uint) (*model.Order, error)
	CancelOrder(ctx context.Context, userID uint, orderID uint) error

	GetMyTrades(ctx context.Context, userID uint, limit, offset int) ([]model.Trade, error)
	GetOrderFills(ctx context.Context, userID uint, orderID uint) ([]model.OrderFill, error)

	// PUBLIC (orderbook / history — later engine will power this)
	GetOrderBook(ctx context.Context, symbol string) ([]model.Order, []model.Order, error)
	GetTradeHistory(ctx context.Context, symbol string, limit int) ([]model.Trade, error)

	// admin side
	GetAllOrders(ctx context.Context, limit, offset int) ([]model.Order, error)
	GetAllTrades(ctx context.Context, limit, offset int) ([]model.Trade, error)
	GetOrdersByUserAdmin(ctx context.Context, userID uint, limit, offset int) ([]model.Order, error)
}



type service struct {
	repo Repository
	redis *redis.Client
	eng *engine.Engine
}

func NewService(repo Repository,rdb *redis.Client,eng *engine.Engine) Service {
	return &service{
		repo: repo,
		redis: rdb,
		eng: eng,
		}
}


type CreateOrderReq struct {
	Symbol   string `json:"symbol"`
	Side     string `json:"side"`   // buy / sell
	Type     string `json:"type"`   // market / limit
	Price    int64  `json:"price"`
	Quantity int64  `json:"quantity"`
}


func (s *service) PlaceOrder(ctx context.Context, userID uint, req CreateOrderReq) error {

	if req.Quantity <= 0 {
		return errors.New("invalid quantity")
	}

	if req.Side != "buy" && req.Side != "sell" {
		return errors.New("invalid side")
	}

	if req.Type == "limit" && req.Price <= 0 {
		return errors.New("price required for limit order")
	}

	var price int64

	if req.Type == "market" {
		p, err := s.getMarketPrice(ctx, req.Symbol)
		if err != nil {
			return err
		}
		price = p
	} else {
		if req.Price <= 0 {
			return errors.New("price required for limit order")
		}
		price = req.Price
	}

	order := &model.Order{
		UserID:       userID,
		Symbol:       req.Symbol,
		Side:         req.Side,
		Type:         req.Type,
		Price:        price,
		Quantity:     req.Quantity,
		RemainingQty: req.Quantity,
		Status:       "open",
	}

	err:= s.repo.CreateOrder(ctx, order)
	if err!=nil{
		return err
	}
	s.eng.Submit(order)

	return nil
}

func (s *service) GetMyOrders(ctx context.Context, userID uint, status string, limit, offset int) ([]model.Order, error) {
	return s.repo.GetOrdersByUser(ctx, userID, status, limit, offset)
}


func (s *service) GetOrderByID(ctx context.Context, userID uint, orderID uint) (*model.Order, error) {

	order, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if order.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	return order, nil
}


func (s *service) CancelOrder(ctx context.Context, userID uint, orderID uint) error {

	order, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.UserID != userID {
		return errors.New("unauthorized")
	}

	if order.Status != "open" {
		return errors.New("only open orders can be cancelled")
	}

	order.Status = "cancelled"

	return s.repo.UpdateOrder(ctx, order)
}

// trade
func (s *service) GetMyTrades(ctx context.Context, userID uint, limit, offset int) ([]model.Trade, error) {
	return s.repo.GetTradesByUser(ctx, userID, limit, offset)
}


func (s *service) GetOrderFills(ctx context.Context, userID uint, orderID uint) ([]model.OrderFill, error) {

	order, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if order.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	return s.repo.GetOrderFills(ctx, orderID)
}

// market read only

// 7. ORDERBOOK (TEMP → DB fallback, later in-memory)
func (s *service) GetOrderBook(ctx context.Context, symbol string) ([]model.Order, []model.Order, error) {

	bids, err := s.repo.GetOpenOrdersForMatching(ctx, symbol, "buy")
	if err != nil {
		return nil, nil, err
	}

	asks, err := s.repo.GetOpenOrdersForMatching(ctx, symbol, "sell")
	if err != nil {
		return nil, nil, err
	}

	return bids, asks, nil
}

// 8. TRADE HISTORY (PUBLIC)
func (s *service) GetTradeHistory(ctx context.Context, symbol string, limit int) ([]model.Trade, error) {
	return s.repo.GetTradesBySymbol(ctx, symbol, limit)
}

//admin side


func (s *service) GetAllOrders(ctx context.Context, limit, offset int) ([]model.Order, error) {
	return s.repo.GetAllOrders(ctx, limit, offset)
}

func (s *service) GetAllTrades(ctx context.Context, limit, offset int) ([]model.Trade, error) {
	return s.repo.GetAllTrades(ctx, limit, offset)
}

func (s *service) GetOrdersByUserAdmin(ctx context.Context, userID uint, limit, offset int) ([]model.Order, error) {
	return s.repo.GetOrdersByUser(ctx, userID, "", limit, offset)
}



func (s *service) getMarketPrice(ctx context.Context, symbol string) (int64, error) {

	key := "market:price:" + symbol

	val, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	price, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, err
	}

	return price, nil
}