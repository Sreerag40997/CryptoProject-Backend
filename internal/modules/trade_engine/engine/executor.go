package engine

import (
	"context"

	"cryptox/internal/modules/trade_engine/model"
)

type WalletService interface {
	DebitINR(ctx context.Context, userID uint, amount int64) error
	CreditINR(ctx context.Context, userID uint, amount int64) error

	DebitCrypto(ctx context.Context, userID uint, symbol string, qty int64) error
	CreditCrypto(ctx context.Context, userID uint, symbol string, qty int64) error
}

type Executor struct {
	repo   Repository
	wallet WalletService
}

func NewExecutor(repo Repository, wallet WalletService) *Executor {
	return &Executor{
		repo:   repo,
		wallet: wallet,
	}
}

func (e *Executor) Execute(
	ctx context.Context,
	buy *model.Order,
	sell *model.Order,
	qty int64,
	price int64,
) error {

	// create trade
	trade := &model.Trade{
		Symbol:      buy.Symbol,
		BuyOrderID:  buy.ID,
		SellOrderID: sell.ID,
		BuyerID:     buy.UserID,
		SellerID:    sell.UserID,
		Price:       price,
		Quantity:    qty,
	}

	if err := e.repo.CreateTrade(ctx, trade); err != nil {
		return err
	}

	// update buy order
	buy.FilledQty += qty
	buy.RemainingQty -= qty

	// update sell order
	sell.FilledQty += qty
	sell.RemainingQty -= qty

	updateStatus(buy)
	updateStatus(sell)

	// save orders
	if err := e.repo.UpdateOrder(ctx, buy); err != nil {
		return err
	}

	if err := e.repo.UpdateOrder(ctx, sell); err != nil {
		return err
	}

	// create fill for buy order
	if err := e.repo.CreateOrderFill(ctx, &model.OrderFill{
		OrderID:  buy.ID,
		TradeID:  trade.ID,
		Price:    price,
		Quantity: qty,
	}); err != nil {
		return err
	}

	// create fill for sell order
	if err := e.repo.CreateOrderFill(ctx, &model.OrderFill{
		OrderID:  sell.ID,
		TradeID:  trade.ID,
		Price:    price,
		Quantity: qty,
	}); err != nil {
		return err
	}

	// wallet settlement
	if err := e.updateWallets(ctx, buy, sell, qty, price); err != nil {
		return err
	}

	return nil
}

func updateStatus(o *model.Order) {

	if o.RemainingQty == 0 {
		o.Status = "filled"
		return
	}

	if o.FilledQty > 0 {
		o.Status = "partial"
		return
	}

	o.Status = "open"
}

func (e *Executor) updateWallets(
	ctx context.Context,
	buy *model.Order,
	sell *model.Order,
	qty int64,
	price int64,
) error {

	amount := qty * price

	// buyer
	if err := e.wallet.DebitINR(ctx, buy.UserID, amount); err != nil {
		return err
	}

	if err := e.wallet.CreditCrypto(ctx, buy.UserID, buy.Symbol, qty); err != nil {
		return err
	}

	// seller
	if err := e.wallet.DebitCrypto(ctx, sell.UserID, sell.Symbol, qty); err != nil {
		return err
	}

	if err := e.wallet.CreditINR(ctx, sell.UserID, amount); err != nil {
		return err
	}

	return nil
}