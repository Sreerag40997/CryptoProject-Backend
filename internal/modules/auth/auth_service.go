package auth

import (
	"context"
	"cryptox/packages/utils"
	"errors"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ctx = context.Background()
)

type AuthService struct {
	repo Repository
	redis *redis.Client
	jwtSecret string
}

func NewAuthService(repo Repository, redis *redis.Client, jwtSecret string) *AuthService {
	return &AuthService{repo: repo, redis: redis, jwtSecret: jwtSecret}
}

// Register func 
func (s *AuthService) Register(req *UserRegisterRequest) (interface{}, error) {

	var existing User
	if err := s.repo.FindOne(&existing, "email = ?", req.Email); err != nil {
		return nil, errors.New("email already exists")
	}

	hashed, err := utils.Hashing(req.Password)
	if err != nil {
		return nil, err
	}

	user := &User{
		Name: req.Name,
		Email: req.Email,
		Password: hashed,
	}

	err = s.repo.Create(user)

	userResponse := &UserRegisterResponse{
		Name: user.Name,
		Email: user.Email,
	}

	return userResponse, err
}

// Login Func
func (s *AuthService) Login(data *UserLoginRequest) (interface{}, string, string, error) {

	var user User
	err := s.repo.FindOne(&user, "email = ?", data.Email)
	if err != nil {
		return nil, "", "",errors.New("user not found")
	}

	if err := utils.Comparepassword(user.Password, data.Password); err != nil {
		return nil, "", "", errors.New("invalid password")
	}

	access, err := utils.GenerateAccess(user.ID, user.Role, s.jwtSecret)
	refresh, err := utils.GenerateRefresh(user.ID, user.Role, s.jwtSecret)
	if err != nil {
		return nil, "", "", errors.New("Token Generating Failed")
	}

	//store in redis
	s.redis.Set(ctx , "refresh"+refresh, user.ID, 7*24*time.Hour)

	loginResponse := &UserLoginResponse{
		ID: user.ID,
		Name: user.Name,
		Email: user.Email,
		Role: user.Role,
		KYCStatus: user.KYCStatus,
		IsVerified: user.IsVerified,
		IsBlocked: user.IsBlocked,
		ProfilePicURL: user.ProfilePicURL,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return loginResponse, access, refresh, nil
}

// Logout
func (s *AuthService) Logout(access, refresh string) {

	s.redis.Del(ctx, "refresh"+refresh)
	s.redis.Set(ctx, "blacklist:"+access, "1", 15*time.Minute)
}



// Refresh Rotation
func (s *AuthService) Refresh(oldRefresh string) (string, string, error) {

	id, err := s.redis.Get(ctx, "refresh"+oldRefresh).Uint64()
	if err != nil {
		return "", "", errors.New("invalid or expired refresh token")
	}

	var user User
	err = s.repo.FindOne(&user, "id = ?", id)
	if err != nil {
		return "", "", errors.New("refresh user notfound")
	}

	s.redis.Del(ctx, "refresh"+oldRefresh)

	newAccess, err := utils.GenerateAccess(user.ID, user.Role, s.jwtSecret)
	newRefresh, err := utils.GenerateRefresh(user.ID, user.Role, s.jwtSecret)
	if err != nil {
		return "", "", errors.New("newAccess or newRefresh token generating failed")
	}

	s.redis.Set(ctx, "refresh"+newRefresh, user.ID, 7*24*time.Hour)

	return newAccess, newRefresh, nil
}

// -> OTP generate busniess logics
func (s *AuthService) SentOtpService(email string) (OTP string, err error) {

	isOk, err := utils.RateLimitOTP(email)

	if err != nil {
		return "", err
	}

	if !isOk {
		return "", errors.New("Request limit exceeded, wait for 10 min")
	}

	RandOTP := utils.GenerateOTP()
	log.Print("OTP:", RandOTP)

	if err := utils.SentOTPEmail(email, RandOTP); err != nil {
		return "", err
	}

	if err := utils.SaveOTP(email, RandOTP); err != nil {
		return "", errors.Join(errors.New("failed to save the OTP"), err)
	}

	return RandOTP, nil
}

// -> Verify the email logic
func (s *AuthService) VerifyOTP(email, OTP string) error {

	storedOtp, err := utils.GetOTP(email)

	if err != nil {
		return errors.Join(errors.New("OTP mismatched"), err)
	}

	if storedOtp != OTP {
		return errors.New("invalid otp")
	}

	utils.DeleteOTP(email)
	return nil
}