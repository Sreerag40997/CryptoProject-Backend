package cashwallet

import "time"

type Wallet struct {
	ID uint `gorm:"primaryKey"`

	UserID uint `gorm:"uniqueIndex;not null"`

	Type string `gorm:"type:varchar(20);default:'user'"` // user | platform

	WalletID string `gorm:"uniqueIndex;not null"` // public ID 

	Balance int64 `gorm:"not null;default:0"` // stored in paise

	PinHash string `gorm:"type:varchar(255)"`

	Currency string `gorm:"type:varchar(10);default:'INR'"`

	Status string `gorm:"type:varchar(20);default:'active'"` 
	// active | frozen | blocked

	CreatedAt time.Time
	UpdatedAt time.Time
}

type WalletTransaction struct {
	ID uint `gorm:"primaryKey"`

	UserID *uint `gorm:"index"`

	WalletID string `gorm:"index;not null"` // matches Wallet.WalletID

	TxnID string `gorm:"uniqueIndex;not null"` // public transaction ID

	Type string `gorm:"type:varchar(20);not null"`
	// credit | debit

	Source string `gorm:"type:varchar(30);not null"`
	// deposit | withdraw | trade | card | refund | admin

	Amount int64 `gorm:"not null"` // in paise

	BalanceAfter int64 `gorm:"not null"` // snapshot after txn

	Status string `gorm:"type:varchar(20);default:'success'"`
	// pending | success | failed

	Reference string `gorm:"type:varchar(255)"`
	// razorpay_payment_id / order_id / trade_id

	Description string `gorm:"type:text"`

	CreatedAt time.Time
}