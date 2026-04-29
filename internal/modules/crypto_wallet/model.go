package cryptowallet

import "time"

type CryptoWallet struct {
	ID      uint  `gorm:"primaryKey"`
	UserID  uint  `gorm:"uniqueIndex:idx_user_asset;not null"`

	AssetID uint  `gorm:"uniqueIndex:idx_user_asset;not null"` //this is the actual colum of the foreinkey
	Asset     CryptoAsset `gorm:"foreignKey:AssetID"`// this is the referal

	Balance int64 `gorm:"not null;default:0"`
	Locked  int64 `gorm:"not null;default:0"` // frozen when sell order is placed

	Status         string `gorm:"default:'active'"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CryptoAsset struct {
	ID uint `gorm:"primaryKey"`

	Symbol         string `gorm:"uniqueIndex;not null"` // BTC, ETH
	Name           string
	Precision      int
	Status         string `gorm:"default:'active'"`
	MinTradeAmount int64  // optional (important for trade engine)

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CryptoWalletTransaction struct {
	ID uint `gorm:"primaryKey"`

	UserID  uint `gorm:"index"`
	AssetID uint `gorm:"index"`

	TxnID string `gorm:"uniqueIndex"`

	Type string // credit | debit | lock | unlock

	Amount int64

	BalanceAfter int64
	LockedAfter  int64

	Reference string // trade_id / admin 
	Status  string `gorm:"default:'success'"`

	Description string

	CreatedAt time.Time
}
