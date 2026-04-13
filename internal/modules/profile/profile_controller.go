package profile

import (
	"cryptox/packages/utils"

	"github.com/gofiber/fiber/v2"
)

type ProfileController struct {
	service ProfileService
}

func NewProfileController(s *ProfileService) *ProfileController {
	return &ProfileController{service: *s}
}

// Profile
func (s *ProfileController) Profile(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	user, err := s.service.Profile(userID)
	if err != nil {
		return utils.Error(c, 500, "User Not Found", err)
	}

	return utils.Success(c, 200, "Successfull", user)
}

// Edit Profile Func
func (s *ProfileController) EditProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var req EditProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "invalied input", err)
	}

	if err := utils.Validator.Struct(&req); err != nil {
		return utils.Error(c, 400, "Input Validation Failed", err)
	}

	editedProfile, err := s.service.EditProfile(req, userID)
	if err != nil {
		return utils.Error(c, 500, "Profile editing failed", err)
	}

	return utils.Success(c, 200, "Profile Edited Successfully", editedProfile)
}

// ChangePassWord
func (s *ProfileController) ChangePassWord(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var req ChangePassWordReq
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "invalied input", err)
	}

	if err := utils.Validator.Struct(&req); err != nil {
		return utils.Error(c, 400, "Input Validation Failed", err)
	}

	if err := s.service.ChangePassowrd(req, userID); err != nil {
		return utils.Error(c, 500, "error", err.Error())
	}

	return utils.Success(c, 200, "Password Changned Successfully", nil)
}

// Delete Account
func (s *ProfileController) DeleteAccount(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	if err := s.service.DeleteAccount(userID); err != nil {
		return utils.Error(c, 500, "error", err)
	}

	return utils.Success(c, 200, "Deleted Successfully", nil)
}