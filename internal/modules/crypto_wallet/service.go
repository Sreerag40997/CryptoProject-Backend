package cryptowallet

import (
	"context"
	"errors"
)

type Service interface {

	//user sidde
	CreateWallet(ctx context.Context, userID uint, symbol string) error
	GetMyWallets(ctx context.Context, userID uint) ([]CryptoWallet, error)
	GetWallet(ctx context.Context, userID uint, symbol string) (*CryptoWallet, error)
	GetSummary(ctx context.Context, userID uint) (int64, error)
	GetTransactions(ctx context.Context, userID uint, symbol string, limit, offset int) ([]CryptoWalletTransaction, error)
	GetLocks(ctx context.Context, userID uint) ([]CryptoWallet, error)

	// admin side
	GetAllWalletsAdmin(ctx context.Context) ([]CryptoWallet, error)
	GetUserWalletsAdmin(ctx context.Context, userID uint) ([]CryptoWallet, error)
	GetUserWalletBySymbolAdmin(ctx context.Context, userID uint, symbol string) (*CryptoWallet, error)

	AdminCredit(ctx context.Context, userID uint, symbol string, amount int64) error
	AdminDebit(ctx context.Context, userID uint, symbol string, amount int64) error

	FreezeWallet(ctx context.Context, userID uint) error
	UnfreezeWallet(ctx context.Context, userID uint) error

	GetAllTransactionsAdmin(ctx context.Context, limit, offset int) ([]CryptoWalletTransaction, error)
	GetUserTransactionsAdmin(ctx context.Context, userID uint, limit, offset int) ([]CryptoWalletTransaction, error)

	///////////asset side///////////
	CreateAsset(ctx context.Context, symbol, name string, precision int) error
	UpdateAsset(ctx context.Context, assetID uint, name string, precision int) error
	UpdateAssetStatus(ctx context.Context, assetID uint, status string) error
	GetAssets(ctx context.Context) ([]CryptoAsset, error)

	///////// trade related//////////
	AddBalance(ctx context.Context, userID uint, symbol string, amount int64) error
	DeductBalance(ctx context.Context, userID uint, symbol string, amount int64) error
	LockBalance(ctx context.Context, userID uint, symbol string, amount int64) error
	UnlockBalance(ctx context.Context, userID uint, symbol string, amount int64) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

/////////// help funcitons ///////////////

func (s *service) getAsset(ctx context.Context, symbol string) (*CryptoAsset, error) {
	asset, err := s.repo.GetAssetBySymbol(ctx, symbol)
	if err != nil {
		return nil, errors.New("asset not found")
	}
	if asset.Status != "active" {
		return nil, errors.New("asset not active")
	}
	return asset, nil
}

func (s *service) getOrCreateWallet(ctx context.Context, userID uint, assetID uint) error {
	_, err := s.repo.GetWallet(ctx, userID, assetID)
	if err == nil {
		return nil
	}

	wallet := &CryptoWallet{
		UserID:  userID,
		AssetID: assetID,
	}

	return s.repo.CreateWallet(ctx, wallet)
}

//////////// user side /////////////////

func (s *service) CreateWallet(ctx context.Context, userID uint, symbol string) error {
	asset, err := s.getAsset(ctx, symbol)
	if err != nil {
		return err
	}
	return s.getOrCreateWallet(ctx, userID, asset.ID)
}

func (s *service) GetMyWallets(ctx context.Context, userID uint) ([]CryptoWallet, error) {
	return s.repo.GetAllWallets(ctx, userID)
}

func (s *service) GetWallet(ctx context.Context, userID uint, symbol string) (*CryptoWallet, error) {
	asset, err := s.getAsset(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return s.repo.GetWallet(ctx, userID, asset.ID)
}

func (s *service) GetSummary(ctx context.Context, userID uint) (int64, error) {
	wallets, err := s.repo.GetAllWallets(ctx, userID)
	if err != nil {
		return 0, err
	}

	var total int64
	for _, w := range wallets {
		total += w.Balance
	}

	return total, nil
}

func (s *service) GetTransactions(
	ctx context.Context,
	userID uint,
	symbol string,
	limit, offset int,
) ([]CryptoWalletTransaction, error) {

	var assetID *uint

	// if symbol provided → convert to assetID
	if symbol != "" {
		asset, err := s.repo.GetAssetBySymbol(ctx, symbol)
		if err != nil {
			return nil, errors.New("asset not found")
		}
		assetID = &asset.ID
	}

	return s.repo.GetTransactions(ctx, userID, assetID, limit, offset)
}

func (s *service) GetLocks(ctx context.Context, userID uint) ([]CryptoWallet, error) {
	wallets, err := s.repo.GetAllWallets(ctx, userID)
	if err != nil {
		return nil, err
	}

	var locked []CryptoWallet
	for _, w := range wallets {
		if w.Locked > 0 {
			locked = append(locked, w)
		}
	}
	return locked, nil
}

//////////// admin side /////////////

func (s *service) GetAllWalletsAdmin(ctx context.Context) ([]CryptoWallet, error) {
	return s.repo.GetAllWalletsAdmin(ctx)
}

func (s *service) GetUserWalletsAdmin(ctx context.Context, userID uint) ([]CryptoWallet, error) {
	return s.repo.GetAllWallets(ctx, userID)
}

func (s *service) GetUserWalletBySymbolAdmin(ctx context.Context, userID uint, symbol string) (*CryptoWallet, error) {
	return s.GetWallet(ctx, userID, symbol)
}

func (s *service) AdminCredit(ctx context.Context, userID uint, symbol string, amount int64) error {
	return s.AddBalance(ctx, userID, symbol, amount)
}

func (s *service) AdminDebit(ctx context.Context, userID uint, symbol string, amount int64) error {
	return s.DeductBalance(ctx, userID, symbol, amount)
}

func (s *service) FreezeWallet(ctx context.Context, userID uint) error {
	return s.repo.UpdateWalletStatus(ctx, userID, "frozen")
}

func (s *service) UnfreezeWallet(ctx context.Context, userID uint) error {
	return s.repo.UpdateWalletStatus(ctx, userID, "active")
}

func (s *service) GetAllTransactionsAdmin(ctx context.Context, limit, offset int) ([]CryptoWalletTransaction, error) {
	return s.repo.GetAllTransactions(ctx, limit, offset)
}

func (s *service) GetUserTransactionsAdmin(ctx context.Context, userID uint, limit, offset int) ([]CryptoWalletTransaction, error) {
	return s.repo.GetTransactions(ctx, userID, nil, limit, offset)
}

//////////// asset management /////////////

func (s *service) CreateAsset(ctx context.Context, symbol, name string, precision int) error {

	if symbol == "" {
		return errors.New("symbol required")
	}

	return s.repo.CreateAsset(ctx, &CryptoAsset{
		Symbol:    symbol,
		Name:      name,
		Precision: precision,
		Status:    "active",
	})
}

func (s *service) UpdateAsset(ctx context.Context, assetID uint, name string, precision int) error {
	return s.repo.UpdateAsset(ctx, &CryptoAsset{
		ID:        assetID,
		Name:      name,
		Precision: precision,
	})
}

func (s *service) UpdateAssetStatus(ctx context.Context, assetID uint, status string) error {
	return s.repo.UpdateAssetStatus(ctx, assetID, status)
}

func (s *service) GetAssets(ctx context.Context) ([]CryptoAsset, error) {
	return s.repo.GetAllAssets(ctx)
}

/////////// trade engine ///////////////

func (s *service) AddBalance(ctx context.Context, userID uint, symbol string, amount int64) error {
	if amount <= 0 {
		return errors.New("invalid amount")
	}

	asset, err := s.getAsset(ctx, symbol)
	if err != nil {
		return err
	}

	if err := s.getOrCreateWallet(ctx, userID, asset.ID); err != nil {
		return err
	}

	return s.repo.AddBalance(ctx, userID, asset.ID, amount)
}

func (s *service) DeductBalance(ctx context.Context, userID uint, symbol string, amount int64) error {
	if amount <= 0 {
		return errors.New("invalid amount")
	}

	asset, err := s.getAsset(ctx, symbol)
	if err != nil {
		return err
	}

	return s.repo.DeductBalance(ctx, userID, asset.ID, amount)
}

func (s *service) LockBalance(ctx context.Context, userID uint, symbol string, amount int64) error {
	asset, err := s.getAsset(ctx, symbol)
	if err != nil {
		return err
	}

	return s.repo.LockBalance(ctx, userID, asset.ID, amount)
}

func (s *service) UnlockBalance(ctx context.Context, userID uint, symbol string, amount int64) error {
	asset, err := s.getAsset(ctx, symbol)
	if err != nil {
		return err
	}

	return s.repo.UnlockBalance(ctx, userID, asset.ID, amount)
}