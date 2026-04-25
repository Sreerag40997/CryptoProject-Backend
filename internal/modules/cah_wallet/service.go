package cashwallet

import (
	"context"
	"cryptox/internal/modules/payment"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	CreateWallet(ctx context.Context, userID uint) error
	// Wallet
	GetMyWallet(ctx context.Context, userID uint) (*Wallet, error)
	GetBalance(ctx context.Context, userID uint) (int64, error)

	// PIN
	SetPin(ctx context.Context, userID uint, pin string) error
	ChangePin(ctx context.Context, userID uint, oldPin, newPin string) error

	// Transactions
	GetTransactions(ctx context.Context, userID uint, limit, offset int) ([]WalletTransaction, error)

	// Money
	Withdraw(ctx context.Context, userID uint, amount int64, pin string) error
	CreateDepositOrder(ctx context.Context, userID uint, amount int64) (string, error)

	// Admin
	AdminGetWallet(ctx context.Context, userID uint) (*Wallet, error)
	AdminBlockWallet(ctx context.Context, userID uint) error
	AdminFreezeWallet(ctx context.Context, userID uint) error
	AdminUnblockWallet(ctx context.Context, userID uint) error
	AdminCredit(ctx context.Context, userID uint, amount int64) error
	AdminDebit(ctx context.Context, userID uint, amount int64) error 

	HandleDepositSuccess(ctx context.Context, userID uint, amount int64, paymentID string) error
	VerifyPayment(orderID, paymentID, signature string) bool
}

type service struct {
	repo Repository
	payment payment.Service
}

func NewService(repo Repository,payment payment.Service) Service {
	return &service{
		repo: repo,
		payment: payment,
	}
}

func generateWalletID() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%014d", rand.Int63n(1e14))
}

func (s *service) CreateWallet(ctx context.Context, userID uint) error {

	// 1. Check if wallet exists
	existingwallet, err := s.repo.GetByUserID(ctx, userID)
	if err == nil && existingwallet != nil {
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

func (s *service) GetMyWallet(ctx context.Context, userID uint) (*Wallet, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *service) GetBalance(ctx context.Context, userID uint) (int64, error) {

	wallet, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}

	return wallet.Balance, nil
}

func (s *service) SetPin(ctx context.Context, userID uint, pin string) error {

	// basic validation
	if len(pin) != 4 {
		return errors.New("PIN must be 4 digits")
	}

	// get wallet
	wallet, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// check if already set
	if wallet.PinHash != "" {
		return errors.New("PIN already set")
	}

	// hash PIN
	hash, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// save
	return s.repo.UpdatePin(ctx, userID, string(hash))
}

func (s *service) verifyPin(wallet *Wallet, pin string) error {

	if wallet.PinHash == "" {
		return errors.New("PIN not set")
	}
	fmt.Println(pin)

	err := bcrypt.CompareHashAndPassword([]byte(wallet.PinHash), []byte(pin))
	if err != nil {
		return errors.New("invalid PIN")
	}

	return nil
}

func (s *service) ChangePin(ctx context.Context, userID uint, oldPin, newPin string) error {

	if len(newPin) != 4 {
		return errors.New("PIN must be 4 digits")
	}

	wallet, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// verify old PIN
	if err := s.verifyPin(wallet, oldPin); err != nil {
		return err
	}

	// hash new PIN
	hash, err := bcrypt.GenerateFromPassword([]byte(newPin), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.UpdatePin(ctx, userID, string(hash))
}

func (s *service) GetTransactions(ctx context.Context, userID uint, limit, offset int) ([]WalletTransaction, error) {
	return s.repo.GetTransactionsByUser(ctx, userID, limit, offset)
}

func (s *service) Withdraw(ctx context.Context, userID uint, amount int64, pin string) error {

	// 1. Validate amount
	if amount <= 0 {
		return errors.New("invalid amount")
	}

	// 2. Get wallet
	wallet, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// 3. Check wallet status
	if wallet.Status != "active" {
		return errors.New("wallet is not active")
	}

	// 4. Verify PIN
	if err := s.verifyPin(wallet, pin); err != nil {
		return err
	}

	// 5. Check balance
	if wallet.Balance < amount {
		return errors.New("insufficient balance")
	}

	// 6. Prepare transaction (PENDING)
	txn := &WalletTransaction{
		TxnID:       generateTxnID(),
		UserID:      &userID,
		WalletID:    wallet.WalletID,
		Type:        "debit",
		Source:      "withdraw",
		Amount:      amount,
		Status:      "pending", // important
		Description: "Withdrawal initiated",
	}

	// 7. DEBIT FIRST (critical)
	if err := s.repo.Debit(ctx, userID, amount, txn); err != nil {
		return err
	}

	// 8. CALL PAYOUT (mock or real)
	payoutID, err := s.payment.CreatePayout(
		userID,
		amount,
		"User Name",   // TODO: fetch from KYC later
		"HDFC0001234", // TODO: fetch from DB
		"1234567890",  // TODO: fetch from DB
	)

	if err != nil {

		// 🔁 ROLLBACK (very important)
		rollbackTxn := &WalletTransaction{
			TxnID:       generateTxnID(),
			UserID:      &userID,
			WalletID:    wallet.WalletID,
			Type:        "credit",
			Source:      "withdraw_failed",
			Amount:      amount,
			Status:      "success",
			Description: "Rollback failed withdrawal",
		}

		_ = s.repo.Credit(ctx, userID, amount, rollbackTxn)

		return err
	}

	// 9. Mark transaction as success
	txn.Status = "success"
	txn.Reference = payoutID

	// update transaction record
	return s.repo.UpdateTransaction(ctx, txn)
}

func (s *service) CreateDepositOrder(ctx context.Context, userID uint, amount int64) (string, error) {

	if amount <= 0 {
		return "", errors.New("invalid amount")
	}

	// check wallet exists
	_, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return "", err
	}

	// call payment module (Razorpay)
	orderID, err := s.payment.CreateOrder(amount, userID)
	if err != nil {
		return "", err
	}

	return orderID, nil
}

func (s *service) AdminGetWallet(ctx context.Context, userID uint) (*Wallet, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *service) AdminBlockWallet(ctx context.Context, userID uint) error {
	return s.repo.UpdateStatus(ctx, userID, "blocked")
}

func (s *service) AdminFreezeWallet(ctx context.Context, userID uint) error {
	return s.repo.UpdateStatus(ctx, userID, "frozen")
}

func (s *service) AdminUnblockWallet(ctx context.Context, userID uint) error {
	return s.repo.UpdateStatus(ctx, userID, "active")
}

func (s *service) AdminCredit(ctx context.Context, userID uint, amount int64) error {

	wallet, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if wallet.Status != "active" {
		return errors.New("wallet not active")
	}
	txn := &WalletTransaction{
		TxnID:       generateTxnID(),
		Source:      "admin",
		Status:      "success",
		Description: "Admin credit",
	}

	return s.repo.Credit(ctx, userID, amount, txn)
}

func (s *service) AdminDebit(ctx context.Context, userID uint, amount int64) error {

	wallet, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if wallet.Status != "active" {
		return errors.New("wallet not active")
	}

	txn := &WalletTransaction{
		TxnID:       generateTxnID(),
		Source:      "admin",
		Status:      "success",
		Description: "Admin debit",
	}

	return s.repo.Debit(ctx, userID, amount, txn)
}

func generateTxnID() string {
	return fmt.Sprintf("TXN_%d", time.Now().UnixNano())
}


func (s *service) HandleDepositSuccess(ctx context.Context, userID uint, amount int64, paymentID string) error {

	//  prevent duplicate credit
	existing, _ := s.repo.GetTransactionByReference(ctx, paymentID)
	if existing != nil {
		return nil
	}

	txn := &WalletTransaction{
		TxnID:      generateTxnID(),
		Source:     "deposit",
		Status:     "success",
		Reference:  paymentID,
		Description:"Added via Razorpay",
	}

	return s.repo.Credit(ctx, userID, amount, txn)
}

func (s *service) VerifyPayment(orderID, paymentID, signature string) bool {
	return s.payment.VerifySignature(orderID, paymentID, signature)
}

func (s *service) CreditPlatform(ctx context.Context, amount int64, source string) error {

	txn := &WalletTransaction{
		TxnID:       generateTxnID(),
		UserID:      nil,
		Type:        "credit",
		Source:      source,
		Amount:      amount,
		Description: "platform earning",
	}

	return s.repo.CreditPlatform(ctx, amount, txn)
}