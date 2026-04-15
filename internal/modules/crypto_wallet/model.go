package cryptowallet

import "time"

type CryptoWallet struct {
	ID     uint
	UserID uint

	CreatedAt time.Time
}


type CryptoAsset struct {
	ID       uint
	UserID   uint

	Symbol   string  // BTC, ETH
	Balance  float64

	UpdatedAt time.Time
}