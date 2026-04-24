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

	// Core
	Credit(ctx context.Context, userID uint, amount int64, txn *WalletTransaction) error
	Debit(ctx context.Context, userID uint, amount int64, txn *WalletTransaction) error

	// Admin
	GetAllWallets(ctx context.Context, limit, offset int) ([]Wallet, error)
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

	//  start DB transaction do all safely or fail everything
	dbTx := r.db.WithContext(ctx).Begin()
	if dbTx.Error != nil {
		return dbTx.Error
	}

	var wallet Wallet

	// lock the wallet row 
	err := dbTx.Clauses(clause.Locking{Strength: "UPDATE"}).//lock wallet one at a time
		Where("user_id = ?", userID).
		First(&wallet).Error

	if err != nil {
		dbTx.Rollback()
		return err
	}

	// calculate new balance
	newBalance := wallet.Balance + amount

	// update wallet balance
	err = dbTx.Model(&Wallet{}).
		Where("id = ?", wallet.ID).
		Update("balance", newBalance).Error

	if err != nil {
		dbTx.Rollback()
		return err
	}

	// fill transaction details
	txn.UserID = userID
	txn.WalletID = wallet.WalletID
	txn.Type = "credit"
	txn.Amount = amount
	txn.BalanceAfter = newBalance

	// insert transaction record
	err = dbTx.Create(txn).Error
	if err != nil {
		dbTx.Rollback()
		return err
	}

	// commit everything
	return dbTx.Commit().Error
}

func (r *repo) Debit(ctx context.Context, userID uint, amount int64, txn *WalletTransaction) error {

	// safety check
	if amount <= 0 {
		return errors.New("invalid amount")
	}

	// start DB transaction
	dbTx := r.db.WithContext(ctx).Begin()
	if dbTx.Error != nil {
		return dbTx.Error
	}

	var wallet Wallet

	// lock wallet row
	err := dbTx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ?", userID).
		First(&wallet).Error

	if err != nil {
		dbTx.Rollback()
		return err
	}

	// check balance
	if wallet.Balance < amount {
		dbTx.Rollback()
		return errors.New("insufficient balance")
	}

	//  new balance
	newBalance := wallet.Balance - amount

	// update wallet
	err = dbTx.Model(&Wallet{}).
		Where("id = ?", wallet.ID).
		Update("balance", newBalance).Error

	if err != nil {
		dbTx.Rollback()
		return err
	}

	//  fill transaction
	txn.UserID = userID
	txn.WalletID = wallet.WalletID
	txn.Type = "debit"
	txn.Amount = amount
	txn.BalanceAfter = newBalance

	// insert transaction
	err = dbTx.Create(txn).Error
	if err != nil {
		dbTx.Rollback()
		return err
	}

	// commit
	return dbTx.Commit().Error
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

func (r *repo) GetAllWallets(ctx context.Context, limit, offset int) ([]Wallet, error) {

	var wallets []Wallet

	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&wallets).Error

	return wallets, err
}