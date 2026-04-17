package ecard

import (
	"context"
	"crypto/rand"
	"cryptox/packages/utils"
	"errors"
	"fmt"
	"math/big"
	"os"
	"time"
)

type Service interface {
	CreateCard(ctx context.Context, userID uint) error
	GetMyCard(ctx context.Context, userID uint) (*CardResponse, error)
	BlockCard(ctx context.Context,userID uint)error
	UnblockCard(ctx context.Context,userID uint)error
	ValidateCard(ctx context.Context,userID uint)error

  AdminGetCard(ctx context.Context, userID uint) (*CardResponse, error)
	AdminBlockCard(ctx context.Context, userID uint) error
	AdminUnblockCard(ctx context.Context, userID uint) error
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

func (s *service)BlockCard(ctx context.Context,userID uint)error{
	card,err:=s.repo.GetByUserID(ctx,userID)
	if err!=nil{
		return err
	}

	if card.Status == "blocked"{
		return nil
	}

	return s.repo.UpdateStatus(ctx,userID,"blocked","")
}

func (s *service)UnblockCard(ctx context.Context,userID uint)error{
	card,err:=s.repo.GetByUserID(ctx,userID)
	if err!=nil{
		return err
	}
	if card.Status == "active"{
		return nil
	}
	if card.BlockedBy == "admin"{
		return errors.New("card blocked by admin,cannot unblock")
	}
	return s.repo.UpdateStatus(ctx,userID,"active","")
}

func (s *service)ValidateCard(ctx context.Context,userID uint)error{
	card,err:=s.repo.GetByUserID(ctx,userID)
	if err!=nil{
		return err
	}
	if card.Status != "active"{
		return errors.New("card is not active")
	}
	return nil
}

func (s *service) AdminBlockCard(ctx context.Context, userID uint) error {

	return s.repo.UpdateStatus(ctx, userID, "blocked", "admin")
}

func (s *service) AdminUnblockCard(ctx context.Context, userID uint) error {

	return s.repo.UpdateStatus(ctx, userID, "active", "")
}

func (s *service) AdminGetCard(ctx context.Context, userID uint) (*CardResponse, error) {

	card, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	maskedCard := "**** **** **** " + card.Last4
	expiry := fmt.Sprintf("%02d/%02d", card.ExpiryMonth, card.ExpiryYear%100)

	return &CardResponse{
		CardNumber: maskedCard,
		Expiry:     expiry,
		Status:     card.Status,
	}, nil
}