package ecard

import (
	"context"
	"crypto/rand"
	"cryptox/packages/utils"
	"fmt"
	"math/big"
	"os"
	"time"
)

type Service interface {
	CreateCard(ctx context.Context, userID uint) error
	GetMyCard(ctx context.Context, userID uint) (*CardResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func randomInt(max int64) int64 {
	n, _ := rand.Int(rand.Reader, big.NewInt(max))
	return n.Int64()
}

func generateCardNumber() string {
	return fmt.Sprintf("%014d", randomInt(1e14))
}

func generateCVV() string {
	return fmt.Sprintf("%03d", randomInt(1000))
}

func (s *service) CreateCard(ctx context.Context, userID uint) error {

	// check existing
	_, err := s.repo.GetByUserID(ctx, userID)
	if err == nil {
		return nil
	}

	cardNumber := generateCardNumber()
	last4 := cardNumber[len(cardNumber)-4:]

	cvv := generateCVV()

	now := time.Now()

	// Load key
	keyStr := os.Getenv("ENCRYPTION_KEY")
	key := []byte(keyStr)

	// Encrypt
	encCard, err := utils.Encrypt(cardNumber, key)
	if err != nil {
		return err
	}

	encCVV, err := utils.Encrypt(cvv, key)
	if err != nil {
		return err
	}

	card := &Card{
		UserID: userID,
		CardNumber: encCard,
		Last4: last4,
		CVV: encCVV,
		ExpiryMonth: int(now.Month()),
		ExpiryYear: now.Year() + 5,
	}

	return s.repo.Create(ctx, card)
}

func (s *service) GetMyCard(ctx context.Context, userID uint) (*CardResponse, error) {

	card, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Mask card number using last4
	maskedCard := "**** **** **** " + card.Last4

	// Format expiry (MM/YY)
	expiry := fmt.Sprintf("%02d/%02d", card.ExpiryMonth, card.ExpiryYear%100)

	return &CardResponse{
		CardNumber: maskedCard,
		Expiry:     expiry,
		Status:     card.Status,
	}, nil
}