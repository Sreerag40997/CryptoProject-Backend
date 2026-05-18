package walletadapter

import (
	"context"
	"strings"

	cashwallet "cryptox/internal/modules/cah_wallet"
	cryptowallet "cryptox/internal/modules/crypto_wallet"
)

type Service struct {
	cash   cashwallet.Service
	crypto cryptowallet.Service
}

func New(cash cashwallet.Service, crypto cryptowallet.Service) *Service {
	return &Service{
		cash:   cash,
		crypto: crypto,
	}
}

// INR (cash wallet)

func (s *Service) DebitINR(ctx context.Context, userID uint, amount int64) error {
	return s.cash.AdminDebit(ctx, userID, amount)
}

func (s *Service) CreditINR(ctx context.Context, userID uint, amount int64) error {
	return s.cash.AdminCredit(ctx, userID, amount)
}

func cleanSymbol(sym string) string {
	return strings.Split(sym, "-")[0]
}

// CRYPTO (crypto wallet)

func (s *Service) DebitCrypto(ctx context.Context, userID uint, symbol string, qty int64) error {
	return s.crypto.DeductBalance(ctx, userID, cleanSymbol(symbol), qty)
}

func (s *Service) CreditCrypto(ctx context.Context, userID uint, symbol string, qty int64) error {
	return s.crypto.AddBalance(ctx, userID, cleanSymbol(symbol), qty)
}