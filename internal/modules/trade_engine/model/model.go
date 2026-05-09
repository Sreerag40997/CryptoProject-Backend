package model

import "time"

type Order struct { //store the trade details
	ID uint `gorm:"primaryKey"`
	UserID uint `gorm:"index;not null"`
	Symbol string `gorm:"index:idx_symbol_side_price;not null"`
	Side string `gorm:"index:idx_symbol_side_price;type:varchar(10);not null"` // buy or sell

	Type string `gorm:"type:varchar(20);not null"` // market , limit , stop_loss , take_profit

	Price int64  `gorm:"index:idx_symbol_side_price"`// required for limit

	Quantity     int64 `gorm:"not null"`  // total qty
	FilledQty    int64 `gorm:"default:0"` // executed qty
	RemainingQty int64 `gorm:"not null"`  // qty left

	Status string `gorm:"index;type:varchar(20);default:'open'"` // open , partial , filled , cancelled

	TimeInForce string `gorm:"type:varchar(10);default:'GTC'"`

	StopPrice   int64 // for stop loss
	TargetPrice int64 // for take profit

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Trade struct {  // execution happeing
	ID uint `gorm:"primaryKey"`

	Symbol string `gorm:"index;not null"`

	BuyOrderID  uint `gorm:"index"`
	SellOrderID uint `gorm:"index"`

	BuyerID  uint
	SellerID uint

	Price    int64 `gorm:"not null"`
	Quantity int64 `gorm:"not null"`

	Fee int64

	CreatedAt time.Time
}

type OrderFill struct {//one piece of an order
	ID uint `gorm:"primaryKey"`

	OrderID uint `gorm:"index"`

	TradeID uint `gorm:"index"`

	Price int64
	Quantity int64

	CreatedAt time.Time
}
