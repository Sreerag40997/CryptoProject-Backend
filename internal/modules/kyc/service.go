package kyc

import (
	"context"
	cashwallet "cryptox/internal/modules/cah_wallet"
	ecard "cryptox/internal/modules/e_card"
	"cryptox/packages/cloudinary"
	"cryptox/packages/utils"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"regexp"
	"sync"
	"time"
)

type Service interface {
	SubmitKYC(ctx context.Context, userID uint, req *SubmitKYCRequest) error
	GetKYCStatus(ctx context.Context, userID uint) (map[string]interface{}, error)
	GetMyKYC(ctx context.Context, userID uint) (*KYCResponse, error) 
	UpdateKYC(ctx context.Context, userID uint, req *UpdateKYCRequest) error

	GetKYCList(ctx context.Context, status string, page int, limit int) ([]KYC, error)
	GetKYCByID(ctx context.Context, id uint) (*KYC, error) 
	UpdateKYCStatus(ctx context.Context, id uint, status string, reason string) error
}

type service struct {
	repo Repository

	cashWalletService cashwallet.Service
	ecardService  ecard.Service
}

func NewService(repo Repository,cashWalletService cashwallet.Service, cardService ecard.Service) Service {
	return &service{
		repo: repo,
		cashWalletService: cashWalletService,
		ecardService: cardService,
	}
}

// submit kyc
func (s *service) SubmitKYC(ctx context.Context, userID uint, req *SubmitKYCRequest) error {
  fmt.Println(req)
	// validate pan
	panRegex := regexp.MustCompile(`^[A-Z]{5}[0-9]{4}[A-Z]$`)
	if !panRegex.MatchString(req.PANNumber) {
		return errors.New("invalid PAN format")
	}

	// validate ifsc
	if len(req.IFSCCode) != 11 {
		return errors.New("invalid IFSC code")
	}


	//  Mask Aadhaar
	maskedAadhaar := "XXXX-XXXX-" + req.AadhaarNumber[len(req.AadhaarNumber)-4:]

	// Mask PAN
	maskedPAN := req.PANNumber[:5] + "****" + req.PANNumber[9:]

	// encrypt account number
	secretKeySr := os.Getenv("ENCRYPTION_KEY")
	key:=[]byte(secretKeySr)

	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
	return errors.New("invalid encryption key length")
  }

	if len(req.AccountNumber) < 8 {
	return errors.New("invalid account number")
  }


	encryptedAccount, err := utils.Encrypt(req.AccountNumber, key)
	if err != nil {
		return err
	}

	last4 := req.AccountNumber[len(req.AccountNumber)-4:]

	// Upload files concurrently

	var wg sync.WaitGroup
	wg.Add(4)

	var aadhaarFrontURL, aadhaarBackURL, panURL, selfieURL string
	var uploadErr error
	var mu sync.Mutex

	upload := func(fileHeader *multipart.FileHeader, result *string) {
		defer wg.Done()

		file, err := fileHeader.Open()
		if err != nil {
			mu.Lock()
			uploadErr = err
			mu.Unlock()
			return
		}
		defer file.Close()

		url, err := cloudinary.UploadFile(file, fileHeader.Filename)
		if err != nil {
			mu.Lock()
			uploadErr = err
			mu.Unlock()
			return
		}

		*result = url
	}

	go upload(req.AadhaarFront, &aadhaarFrontURL)
	go upload(req.AadhaarBack, &aadhaarBackURL)
	go upload(req.PANFile, &panURL)
	go upload(req.Selfie, &selfieURL)

	wg.Wait()

	if uploadErr != nil {
		return uploadErr
	}

	// parse DBO
	dob, err := time.Parse("2006-01-02", req.DOB)
	if err != nil {
		return errors.New("invalid DOB format")
	}

	// create model
	kyc := &KYC{
		UserID: userID,

		FullName: req.FullName,
		DOB:      dob,

		AadhaarNumber: maskedAadhaar,
		PANNumber:     maskedPAN,

		AccountNumber: encryptedAccount,
		AccountLast4: last4,
		IFSC: req.IFSCCode,

		AadhaarFrontURL: aadhaarFrontURL,
		AadhaarBackURL:  aadhaarBackURL,
		PANURL:          panURL,
		SelfieURL:       selfieURL,

		Status: "pending",
	}

	// save to db
	return s.repo.Create(ctx, kyc)
}


// for kyc status
func (s *service) GetKYCStatus(ctx context.Context, userID uint) (map[string]interface{}, error) {

	kyc, err := s.repo.GetByUserID(ctx, userID)

	// No KYC found
	if err != nil {
		return map[string]interface{}{
			"status":  "not_submitted",
			"reason":  nil,
		}, nil
	}

	switch kyc.Status {

	case "pending":
		return map[string]interface{}{
			"status": "pending",
			"reason": nil,
		}, nil

	case "approved":
		return map[string]interface{}{
			"status": "approved",
			"reason": nil,
		}, nil

	case "rejected":
		return map[string]interface{}{
			"status": "rejected",
			"reason": kyc.RejectionReason,
		}, nil
	}

	return nil, nil
}

// get kyc

func (s *service) GetMyKYC(ctx context.Context, userID uint) (*KYCResponse, error) {

	kyc, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// mask ac number
	accountMasked := "****"
	if len(kyc.AccountNumber) >= 4 {
		accountMasked = "****" + kyc.AccountNumber[len(kyc.AccountNumber)-4:]
	}

	return &KYCResponse{
		ID:       kyc.ID,
		FullName: kyc.FullName,
		Status:   kyc.Status,

		AadhaarMasked: kyc.AadhaarNumber,
		PANMasked:     kyc.PANNumber,

		AccountMasked: accountMasked,
		IFSCCode:      kyc.IFSC,

		AadhaarFrontURL: kyc.AadhaarFrontURL,
		AadhaarBackURL:  kyc.AadhaarBackURL,
		PANURL:          kyc.PANURL,
		SelfieURL:       kyc.SelfieURL,
	}, nil
}

// update kyc
func (s *service) UpdateKYC(ctx context.Context, userID uint, req *UpdateKYCRequest) error {

	kyc, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return errors.New("kyc not found")
	}

	// only for rejected
	if kyc.Status != "rejected" {
		return errors.New("kyc cannot be updated")
	}

	// Optional updates
	if req.FullName != "" {
		kyc.FullName = req.FullName
	}

	if req.AadhaarNumber != "" {
		kyc.AadhaarNumber = "XXXX-XXXX-" + req.AadhaarNumber[len(req.AadhaarNumber)-4:]
	}

	if req.PANNumber != "" {
		kyc.PANNumber = req.PANNumber[:5] + "****" + req.PANNumber[9:]
	}

	if req.AccountNumber != "" {
		// encrypt
		encrypted, err := utils.Encrypt(req.AccountNumber, []byte("1234567890123456"))
		if err != nil {
			return err
		}
		kyc.AccountNumber = encrypted
	}

	if req.IFSCCode != "" {
		kyc.IFSC = req.IFSCCode
	}

	// Reset status
	kyc.Status = "pending"
	kyc.RejectionReason = ""

	return s.repo.Update(ctx, kyc)
}

func (s *service) GetKYCList(ctx context.Context, status string, page int, limit int) ([]KYC, error) {

	offset := (page - 1) * limit

	return s.repo.ListWithFilter(ctx, status, limit, offset)
}

func (s *service) GetKYCByID(ctx context.Context, id uint) (*KYC, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) UpdateKYCStatus(ctx context.Context, id uint, status string, reason string) error {

	kyc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.repo.UpdateStatus(ctx, id, status, reason)
	if err != nil {
		return err
	}

	if status == "approved" {
		userID := kyc.UserID

		_ = s.cashWalletService.CreateWallet(ctx, userID)
		_ = s.ecardService.CreateCard(ctx, userID)
	}

	return nil
}