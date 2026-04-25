package cashwallet

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {

	// Wallet
	Create(ctx context.Context, wallet *Wallet) error
	GetByUserID(ctx context.Context, userID uint) (*Wallet, error)
	GetByWalletID(ctx context.Context, walletID string) (*Wallet, error)
	UpdateWallet(ctx context.Context, wallet *Wallet) error
	UpdateStatus(ctx context.Context, userID uint, status string) error
	UpdatePin(ctx context.Context, userID uint, pinHash string) error

	// Transactions
	CreateTransaction(ctx context.Context, txn *WalletTransaction) error
	GetTransactionsByUser(ctx context.Context, userID uint, limit, offset int) ([]WalletTransaction, error)
	GetTransactionByTxnID(ctx context.Context, txnID string) (*WalletTransaction, error)
	GetTransactionByReference(ctx context.Context, ref string) (*WalletTransaction, error)
	UpdateTransaction(ctx context.Context, txn *WalletTransaction) error

	// Core
	Credit(ctx context.Context, userID uint, amount int64, txn *WalletTransaction) error
	Debit(ctx context.Context, userID uint, amount int64, txn *WalletTransaction) error

	// Admin
	GetAllWallets(ctx context.Context, limit, offset int) ([]Wallet, error)
	GetPlatformWallet(ctx context.Context) (*Wallet, error)
	CreditPlatform(ctx context.Context, amount int64, txn *WalletTransaction) error
}
type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{
		db: db,
	}
}


func (r *repo) Create(ctx context.Context, wallet *Wallet) error {
	return r.db.WithContext(ctx).Create(wallet).Error
}

func (r *repo) GetByUserID(ctx context.Context, userID uint) (*Wallet, error) {
	var wallet Wallet

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&wallet).Error

	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (r *repo) Credit(ctx context.Context, userID uint, amount int64, txn *WalletTransaction) error {

	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	var wallet Wallet

	// lock wallet row
	err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ?", userID).
		First(&wallet).Error

	if err != nil {
		return err
	}

	newBalance := wallet.Balance + amount

	// update balance
	err = r.db.WithContext(ctx).
		Model(&Wallet{}).
		Where("id = ?", wallet.ID).
		Update("balance", newBalance).Error

	if err != nil {
		return err
	}

	// fill txn
	txn.UserID = &userID
	txn.WalletID = wallet.WalletID
	txn.Type = "credit"
	txn.Amount = amount
	txn.BalanceAfter = newBalance

	return r.db.WithContext(ctx).Create(txn).Error
}

func (r *repo) Debit(ctx context.Context, userID uint, amount int64, txn *WalletTransaction) error {

	if amount <= 0 {
		return errors.New("invalid amount")
	}

	var wallet Wallet

	// lock wallet
	err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ?", userID).
		First(&wallet).Error

	if err != nil {
		return err
	}

	if wallet.Balance < amount {
		return errors.New("insufficient balance")
	}

	newBalance := wallet.Balance - amount

	// update
	err = r.db.WithContext(ctx).
		Model(&Wallet{}).
		Where("id = ?", wallet.ID).
		Update("balance", newBalance).Error

	if err != nil {
		return err
	}

	// fill txn
	txn.UserID = &userID
	txn.WalletID = wallet.WalletID
	txn.Type = "debit"
	txn.Amount = amount
	txn.BalanceAfter = newBalance

	return r.db.WithContext(ctx).Create(txn).Error
}

func (r *repo) GetByWalletID(ctx context.Context, walletID string) (*Wallet, error) {
	var wallet Wallet
	err := r.db.WithContext(ctx).
		Where("wallet_id = ?", walletID).
		First(&wallet).Error

	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *repo) UpdateWallet(ctx context.Context, wallet *Wallet) error {
	return r.db.WithContext(ctx).Save(wallet).Error
}

func (r *repo) UpdateStatus(ctx context.Context, userID uint, status string) error {
	return r.db.WithContext(ctx).
		Model(&Wallet{}).
		Where("user_id = ?", userID).
		Update("status", status).Error
}

func (r *repo) UpdatePin(ctx context.Context, userID uint, pinHash string) error {
	return r.db.WithContext(ctx).
		Model(&Wallet{}).
		Where("user_id = ?", userID).
		Update("pin_hash", pinHash).Error
}

func (r *repo) CreateTransaction(ctx context.Context, txn *WalletTransaction) error {
	return r.db.WithContext(ctx).Create(txn).Error
}

func (r *repo) GetTransactionsByUser(ctx context.Context, userID uint, limit, offset int) ([]WalletTransaction, error) {

	var txns []WalletTransaction

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&txns).Error

	return txns, err
}

func (r *repo) GetTransactionByTxnID(ctx context.Context, txnID string) (*WalletTransaction, error) {

	var txn WalletTransaction

	err := r.db.WithContext(ctx).
		Where("txn_id = ?", txnID).
		First(&txn).Error

	if err != nil {
		return nil, err
	}

	return &txn, nil
}

func (r *repo) GetTransactionByReference(ctx context.Context, ref string) (*WalletTransaction, error) {

	var txn WalletTransaction

	err := r.db.WithContext(ctx).
		Where("reference = ?", ref).
		First(&txn).Error

	if err != nil {
		return nil, err
	}

	return &txn, nil
}

func (r *repo) UpdateTransaction(ctx context.Context, txn *WalletTransaction) error {
	return r.db.WithContext(ctx).
		Model(&WalletTransaction{}).
		Where("txn_id = ?", txn.TxnID).
		Updates(map[string]interface{}{
			"status":    txn.Status,
			"reference": txn.Reference,
		}).Error
}

func (r *repo) GetAllWallets(ctx context.Context, limit, offset int) ([]Wallet, error) {

	var wallets []Wallet

	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&wallets).Error

	return wallets, err
}

func (r *repo) GetPlatformWallet(ctx context.Context) (*Wallet, error) {
	var wallet Wallet
	err := r.db.WithContext(ctx).
		Where("type = ?", "platform").
		First(&wallet).Error

	return &wallet, err
}

func (r *repo) WithTx(ctx context.Context, fn func(Repository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &repo{db: tx}
		return fn(txRepo)
	})
}

func (r *repo) CreditPlatform(ctx context.Context, amount int64, txn *WalletTransaction) error {

	var wallet Wallet
	if err := r.db.WithContext(ctx).
		Where("type = ?", "platform").
		First(&wallet).Error; err != nil {
		return err
	}

	newBalance := wallet.Balance + amount

	err := r.db.WithContext(ctx).
		Model(&wallet).
		Update("balance", newBalance).Error

	if err != nil {
		return err
	}

	txn.WalletID = wallet.WalletID
	txn.UserID = nil
	txn.Type = "credit"
	txn.Amount = amount
	txn.BalanceAfter = newBalance

	return r.db.WithContext(ctx).Create(txn).Error
}