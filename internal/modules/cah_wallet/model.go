package cashwallet

import "time"

type Wallet struct {
	ID        uint      `gorm:"primaryKey"`

	UserID    uint      `gorm:"uniqueIndex;not null"` 

	WalletID  string    `gorm:"uniqueIndex;not null"` 

	Balance   float64   `gorm:"type:decimal(20,8);default:0"` 

	Currency  string    `gorm:"type:varchar(10);default:'INR'"`

	PinHash   string    `gorm:"type:varchar(255)"` 

	Status    string    `gorm:"type:varchar(20);default:'active'"` // active, frozen, blocked
	

	CreatedAt time.Time
	UpdatedAt time.Time
}

type WalletTransaction struct {
	ID        uint      `gorm:"primaryKey"`

	UserID    uint      `gorm:"index;not null"`
	WalletID  uint      `gorm:"index;not null"`

	TxnID     string    `gorm:"uniqueIndex;not null"` // public txn id

	Type      string    `gorm:"type:varchar(20);not null"` // credit / debit

	Amount    float64   `gorm:"type:decimal(20,8);not null"`

	BalanceAfter float64 `gorm:"type:decimal(20,8)"` // snapshot

	Status    string    `gorm:"type:varchar(20);default:'success'"`// success / failed / pending

	Reference string    `gorm:"type:varchar(255)"` // razorpay_id 

	Description string  `gorm:"type:text"`

	CreatedAt time.Time
}