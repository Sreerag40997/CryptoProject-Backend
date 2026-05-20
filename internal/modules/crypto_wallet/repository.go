package cryptowallet

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {

	// wallet
	CreateWallet(ctx context.Context, wallet *CryptoWallet) error
	GetWallet(ctx context.Context, userID uint, assetID uint) (*CryptoWallet, error)
	GetWalletBySymbol(ctx context.Context, userID uint, symbol string) (*CryptoWallet, error)
	GetAllWallets(ctx context.Context, userID uint) ([]CryptoWallet, error)
	GetAllWalletsAdmin(ctx context.Context) ([]CryptoWallet, error)
	UpdateWalletStatus(ctx context.Context, userID uint, status string) error

	// balance operations
	AddBalance(ctx context.Context, userID uint, assetID uint, amount int64) error
	DeductBalance(ctx context.Context, userID uint, assetID uint, amount int64) error

	// lock system  for orders
	LockBalance(ctx context.Context, userID uint, assetID uint, amount int64) error
	UnlockBalance(ctx context.Context, userID uint, assetID uint, amount int64) error

	// admin
	AdminCredit(ctx context.Context, userID uint, assetID uint, amount int64) error
	AdminDebit(ctx context.Context, userID uint, assetID uint, amount int64) error

	// asset
	CreateAsset(ctx context.Context, asset *CryptoAsset) error
	UpdateAsset(ctx context.Context, asset *CryptoAsset) error
	UpdateAssetStatus(ctx context.Context, assetID uint, status string) error

	GetAssetByID(ctx context.Context, id uint) (*CryptoAsset, error)
	GetAssetBySymbol(ctx context.Context, symbol string) (*CryptoAsset, error)
	GetAllAssets(ctx context.Context) ([]CryptoAsset, error)

	// transactions
CreateTransaction(ctx context.Context, txn *CryptoWalletTransaction) error
GetTransactions(ctx context.Context, userID uint, assetID *uint, limit, offset int) ([]CryptoWalletTransaction, error)
GetAllTransactions(ctx context.Context, limit, offset int) ([]CryptoWalletTransaction, error)

	// transaction support (optional for now)
	WithTx(ctx context.Context, fn func(Repository) error) error
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{db: db}
}

func (r *repo) CreateWallet(ctx context.Context, wallet *CryptoWallet) error {
	return r.db.WithContext(ctx).Create(wallet).Error
}

//////////////// wallet section  ///////////////

//////// user /////////////

//get wallet by asset
func (r *repo) GetWallet(ctx context.Context, userID uint, assetID uint) (*CryptoWallet, error) {

	var wallet CryptoWallet

	err := r.db.WithContext(ctx).
		Where("user_id = ? AND asset_id = ?", userID, assetID).
		First(&wallet).Error

	return &wallet, err
}

// get wallet by symbol
func (r *repo) GetWalletBySymbol(ctx context.Context, userID uint, symbol string) (*CryptoWallet, error) {

	var wallet CryptoWallet

	err := r.db.WithContext(ctx).
		Joins("JOIN crypto_assets ON crypto_assets.id = crypto_wallets.asset_id").
		Where("crypto_wallets.user_id = ? AND crypto_assets.symbol = ?", userID, symbol).
		First(&wallet).Error

	return &wallet, err
}

func (r *repo) GetAllWallets(ctx context.Context, userID uint) ([]CryptoWallet, error) {

	var wallets []CryptoWallet

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&wallets).Error

	return wallets, err
}

func (r *repo) GetAllWalletsAdmin(ctx context.Context) ([]CryptoWallet, error) {

	var wallets []CryptoWallet

	err := r.db.WithContext(ctx).
		Preload("Asset"). // join asset info (BTC, ETH)
		Order("id DESC").
		Find(&wallets).Error

	return wallets, err
}

func (r *repo) UpdateWalletStatus(ctx context.Context, userID uint, status string) error {

	return r.db.WithContext(ctx).
		Model(&CryptoWallet{}).
		Where("user_id = ?", userID).
		Update("status", status).Error
}

/////////// balance operation ////////////

func (r *repo) AddBalance(ctx context.Context, userID uint, assetID uint, amount int64) error {

	var wallet CryptoWallet

	err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ? AND asset_id = ?", userID, assetID).
		First(&wallet).Error

	if err != nil {
		return err
	}

	newBalance := wallet.Balance + amount

	if err := r.db.WithContext(ctx).
		Model(&wallet).
		Update("balance", newBalance).Error; err != nil {
		return err
	}

	//  create transaction
	txn := &CryptoWalletTransaction{
		TxnID:        generateTxnID(),
		UserID:       userID,
		AssetID:      assetID,
		Type:         "credit",
		Amount:       amount,
		BalanceAfter: newBalance,
		LockedAfter:  wallet.Locked,
		Description:  "balance credited",
	}

	return r.CreateTransaction(ctx, txn)
}

func (r *repo) DeductBalance(ctx context.Context, userID uint, assetID uint, amount int64) error {

	var wallet CryptoWallet

	err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ? AND asset_id = ?", userID, assetID).
		First(&wallet).Error

	if err != nil {
		return err
	}

	if wallet.Balance-wallet.Locked < amount {
		return errors.New("insufficient balance")
	}

	newBalance := wallet.Balance - amount

	if err:= r.db.WithContext(ctx).
		Model(&wallet).
		Update("balance", newBalance).Error;err!=nil{
			return err
		}

		txn := &CryptoWalletTransaction{
	TxnID:        generateTxnID(),
	UserID:       userID,
	AssetID:      assetID,
	Type:         "debit",
	Amount:       amount,
	BalanceAfter: newBalance,
	LockedAfter:  wallet.Locked,
	Description:  "balance debited",
}
return r.CreateTransaction(ctx, txn)
}

///// lock system for buying and sell //////////

func (r *repo) LockBalance(ctx context.Context, userID uint, assetID uint, amount int64) error {

	var wallet CryptoWallet

	err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ? AND asset_id = ?", userID, assetID).
		First(&wallet).Error

	if err != nil {
		return err
	}

	available := wallet.Balance - wallet.Locked

	if available < amount {
		return errors.New("insufficient balance")
	}

	wallet.Locked += amount

	if err := r.db.WithContext(ctx).Save(&wallet).Error; err != nil {
		return err
	}

	//  transaction
	txn := &CryptoWalletTransaction{
		TxnID:        generateTxnID(),
		UserID:       userID,
		AssetID:      assetID,
		Type:         "lock",
		Amount:       amount,
		BalanceAfter: wallet.Balance,
		LockedAfter:  wallet.Locked,
		Description:  "balance locked for order",
	}

	return r.CreateTransaction(ctx, txn)
}

func (r *repo) UnlockBalance(ctx context.Context, userID uint, assetID uint, amount int64) error {

	var wallet CryptoWallet

	err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ? AND asset_id = ?", userID, assetID).
		First(&wallet).Error

	if err != nil {
		return err
	}

	if wallet.Locked < amount {
		return errors.New("invalid unlock amount")
	}

	wallet.Locked -= amount

	if err:= r.db.WithContext(ctx).Save(&wallet).Error;err!=nil{
		return err
	}
	txn := &CryptoWalletTransaction{
	TxnID:        generateTxnID(),
	UserID:       userID,
	AssetID:      assetID,
	Type:         "unlock",
	Amount:       amount,
	BalanceAfter: wallet.Balance,
	LockedAfter:  wallet.Locked,
	Description:  "balance unlocked",
}
return r.CreateTransaction(ctx, txn)
}

////////////// admin operation in wallet///////////////////

func (r *repo) AdminCredit(ctx context.Context, userID uint, assetID uint, amount int64) error {
	return r.AddBalance(ctx, userID, assetID, amount)
}

func (r *repo) AdminDebit(ctx context.Context, userID uint, assetID uint, amount int64) error {
	return r.DeductBalance(ctx, userID, assetID, amount)
}

//////////////// asset section ////////////////

func (r *repo) CreateAsset(ctx context.Context, asset *CryptoAsset) error {
	return r.db.WithContext(ctx).Create(asset).Error
}

func (r *repo) UpdateAsset(ctx context.Context, asset *CryptoAsset) error {
	return r.db.WithContext(ctx).
		Model(&CryptoAsset{}).
		Where("id = ?", asset.ID).
		Updates(asset).Error
}

func (r *repo) UpdateAssetStatus(ctx context.Context, assetID uint, status string) error {
	return r.db.WithContext(ctx).
		Model(&CryptoAsset{}).
		Where("id = ?", assetID).
		Update("status", status).Error
}

func (r *repo) GetAssetByID(ctx context.Context, id uint) (*CryptoAsset, error) {
	var asset CryptoAsset

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&asset).Error

	return &asset, err
}


func (r *repo) GetAssetBySymbol(ctx context.Context, symbol string) (*CryptoAsset, error) {

	var asset CryptoAsset

	err := r.db.WithContext(ctx).
		Where("symbol = ?", symbol).
		First(&asset).Error

	return &asset, err
}

func (r *repo) GetAllAssets(ctx context.Context) ([]CryptoAsset, error) {

	var assets []CryptoAsset

	err := r.db.WithContext(ctx).Find(&assets).Error

	return assets, err
}

//////// transaction related /////////////

func (r *repo) CreateTransaction(ctx context.Context, txn *CryptoWalletTransaction) error {
	return r.db.WithContext(ctx).Create(txn).Error
}

func (r *repo) GetTransactions(ctx context.Context, userID uint, assetID *uint, limit, offset int) ([]CryptoWalletTransaction, error) {

	var txns []CryptoWalletTransaction

	query := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("id DESC").
		Limit(limit).
		Offset(offset)

	if assetID != nil {
		query = query.Where("asset_id = ?", *assetID)
	}

	err := query.Find(&txns).Error
	return txns, err
}

func (r *repo) GetAllTransactions(ctx context.Context, limit, offset int) ([]CryptoWalletTransaction, error) {

	var txns []CryptoWalletTransaction

	err := r.db.WithContext(ctx).
		Order("id DESC").
		Limit(limit).
		Offset(offset).
		Find(&txns).Error

	return txns, err
}

func generateTxnID() string {
	return fmt.Sprintf("CTX-%d-%s", time.Now().UnixNano(), uuid.New().String()[:8])
}

// with this run mutiple db operation as one
func (r *repo) WithTx(ctx context.Context, fn func(Repository) error) error {

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		txRepo := &repo{db: tx}

		return fn(txRepo)
	})
}

