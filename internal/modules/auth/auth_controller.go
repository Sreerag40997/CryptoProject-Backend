package auth

import (
	"cryptox/packages/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	Service *AuthService
}

func NewAuthController(s *AuthService) *AuthController {
	return &AuthController{Service: s}
}

// Registration Func
func (s *AuthController) Register(c *fiber.Ctx) error {

	var newUser UserRegisterRequest
	if err := c.BodyParser(&newUser); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid input",
			"err":   err.Error(),
		})
	}

	if err := utils.Validator.Struct(newUser); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid input",
			"err":   err.Error(),
		})
	}

	user, err := s.Service.Register(&newUser)
	if err != nil {
		return utils.Error(c, 400, "Registration failed", err)
	}

	return utils.Success(c, 200, "Registration Successful", user)
}

// Login Func
func (s *AuthController) Login(c *fiber.Ctx) error {

	var req UserLoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": utils.Error(c, 400, "invalid input", nil),
		})
	}

	user, access, refresh, err := s.Service.Login(&req)

	if err != nil {
		return utils.Error(c, 401, "Unauthorized User", err.Error())
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access",
		Value:    access,
		HTTPOnly: true,
		SameSite: "Lax", // Strict in production
		Secure:   false, // true in production (HTTPS)
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 7, // 15 minutes
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refresh",
		Value:    refresh,
		HTTPOnly: true,
		SameSite: "Lax",
		Secure:   false,
		Path:     "/",              // only sent to refresh endpoint
		MaxAge:   60 * 60 * 24 * 7, // 7 days
	})

	return utils.Success(c, 200, "Login Successful", user)
}

// Logout Func
func (s *AuthController) Logout(c *fiber.Ctx) error {

	s.Service.Logout(c.Cookies("access"), c.Cookies("refresh"))

	return utils.Success(c, 200, "Logout Successful", nil)
}

// Refresh func
func (s *AuthController) Refresh(c *fiber.Ctx) error {

	access, refresh, err := s.Service.Refresh(c.Cookies("refresh"))
	if err != nil {
		return utils.Error(c, 401, "Rotation Failed", err)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access",
		Value:    access,
		HTTPOnly: true,
		SameSite: "Lax",
		Secure:   false,
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 7})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh",
		Value:    refresh,
		HTTPOnly: true,
		SameSite: "Lax",
		Secure:   false,
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 7,
	})

	return utils.Success(c, 200, "Rotated", nil)
}

// sent otp func
func (s *AuthController) SendOTP(c *fiber.Ctx) error {

	var OtpEmail struct {
		Email string `json:"email" validate:"required,email"`
	}

	if err := c.BodyParser(&OtpEmail); err != nil {
		return utils.Error(c, 400, "invalid request", nil)
	}

	if OtpEmail.Email == "" {
		return utils.Error(c, 400, "email required", nil)
	}

	otp, err := s.Service.SentOtpService(OtpEmail.Email)
	if err != nil {
		return utils.Error(c, 500, err.Error(), nil)
	}

	return utils.Success(c, 200, "otp sent", otp)
}

// VerifyOTP
func (s *AuthController) VerifyOTP(c *fiber.Ctx) error {

	var VerifyOtp struct {
		Email string `json:"email" validate:"required,email"`
		Otp   string `json:"otp" validate:"required,len=6,numeric"`
	}

	if err := c.BodyParser(&VerifyOtp); err != nil {
		return utils.Error(c, 400, "Invalied Input", err)
	}

	if err := utils.Validator.Struct(&VerifyOtp); err != nil {
		return utils.Error(c, 400, "Input Validation Failed", err)
	}

	if err := s.Service.VerifyOTP(VerifyOtp.Email, VerifyOtp.Otp); err != nil {
		return utils.Error(c, 500, "Email Verification Failed", err)
	}

	return utils.Success(c, 200, "Email Verified Successfully", nil)
}

// Forgot PassWord Func
func (s *AuthController) ForgotPassWordOTP(c *fiber.Ctx) error {

	var email struct {
		Email string `json:"email" validate:"required,email"`
	}

	if err := c.BodyParser(&email); err != nil {
		return utils.Error(c, 400, "Invalied Input", err)
	}

	if err := utils.Validator.Struct(&email); err != nil {
		return utils.Error(c, 400, "Input Validation Failed", err)
	}

	if err := s.Service.ForgotPassWordOTP(email.Email); err != nil {
		return utils.Error(c, 500, "Forgot Password OTP Sending Failed...", err.Error())
	}

	return utils.Success(c, 200, "OTP Sented", nil)
}

func (s *AuthController) ForgotPassWordNewCreation(c *fiber.Ctx) error {

	var req struct {
		Email           string `json:"email" validate:"required,email"`
		NewPassword     string `json:"newpassword" validate:"required"`
		ConfirmPassword string `json:"confirmpassword" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "Invalied Input", err)
	}

	if err := utils.Validator.Struct(&req); err != nil {
		return utils.Error(c, 400, "Input Validation Failed", err)
	}

	err := s.Service.ForgotPassWordNewCreation(req.Email, req.NewPassword, req.ConfirmPassword)
	if err != nil {
		return utils.Error(c, 500, "Password Changing Fauled", err.Error())
	}

	return utils.Success(c, 200, "Password Changed Successfully", nil)
}

//////////////// Admin Functions \\\\\\\\\\\\\\\\

func (s *AuthController) GetAllUsers(c *fiber.Ctx) error {

	users, err := s.Service.GetAllUsers()
	if err != nil {
		return utils.Error(c, 500, "Users Getting Failed", err)
	}

	return utils.Success(c, 200, "success", users)
}

func (s *AuthController) GetByID(c *fiber.Ctx) error {

	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return utils.Error(c, 400, "invaled userID", err)
	}

	user, err2 := s.Service.GetByID(uint(id))
	if err2 != nil {
		return utils.Error(c, 500, "User Not Found", err2)
	}

	return utils.Success(c, 200, "success", user)
}

func (s *AuthController) EditProfile(c *fiber.Ctx) error {

	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return utils.Error(c, 400, "invaled userID", err)
	}

	var req EditProfileReq
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "invaled input", err)
	}

	editedprofile, err1 := s.Service.EditProfile(uint(id), &req)
	if err1 !=nil {
		return utils.Error(c, 500, "updating failed", err1)
	}

	return utils.Success(c, 200, "success", editedprofile)
}

// Block Unblock
func (s *AuthController) BlockUnblock(c *fiber.Ctx) error {

	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return utils.Error(c, 400, "invaled userID", err.Error())
	}

	status, err := s.Service.BlockUnblock(uint(id))
	if err != nil {
		return utils.Error(c, 500, " BlockUnblock failed", err.Error())
	}

	return utils.Success(c, 200, "success", status)
}