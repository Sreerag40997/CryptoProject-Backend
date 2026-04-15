package cashwallet

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type Service interface {
	CreateWallet(ctx context.Context, userID uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func generateWalletID() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%014d", rand.Int63n(1e14))
}

func (s *service) CreateWallet(ctx context.Context, userID uint) error {

	// 1. Check if wallet exists
	_, err := s.repo.GetByUserID(ctx, userID)
	if err == nil {
		// already exists → do nothing
		return nil
	}

	// 2. Generate wallet ID
	walletID := generateWalletID()

	// 3. Create wallet object
	wallet := &Wallet{
		UserID:   userID,
		WalletID: walletID,
		Balance:  0,
		Currency: "INR",
		Status:   "active",
	}

	// 4. Save to DB
	return s.repo.Create(ctx, wallet)
}