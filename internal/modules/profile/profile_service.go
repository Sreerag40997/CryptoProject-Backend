package profile

import (
	"cryptox/packages/utils"
	"errors"
)


type ProfileService struct {
	repo Repository
}

func NewProfileService(repo Repository) *ProfileService {
	return &ProfileService{repo: repo}
}

// Get Profile
func (s *ProfileService) Profile(userID uint) (*UserProfile, error) {

	var user User
	err := s.repo.FindOne(&user, "id", userID)
	if err != nil {
		return nil, err
	}

	userProfile := &UserProfile{
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

	return userProfile, nil
}

//Edit Profile
func (s *ProfileService) EditProfile(req EditProfileRequest, userID uint) (interface{}, error) {

	var user User
	if err := s.repo.FindOne(&user, "id = ?", userID); err != nil {
		return nil, err
	}

	if req.Name != "" {
		field := make(map[string]interface{})
		field["name"] = req.Name
		if err := s.repo.Update(&user, field, "id = ?", userID); err != nil {
			return nil, err
		}
	}

	editedUser := &UserProfile{
		ID: user.ID,
		Name: req.Name,
		Email: user.Email,
		Role: user.Role,
		KYCStatus: user.KYCStatus,
		IsVerified: user.IsVerified,
		IsBlocked: user.IsBlocked,
		ProfilePicURL: user.ProfilePicURL,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	
	return editedUser, nil
}

// ChangePassword
func (s *ProfileService) ChangePassowrd(req ChangePassWordReq, userID uint) error {

	var user User
	if err := s.repo.FindOne(&user, "id = ?", userID); err != nil {
		return err
	}

	if req.OldPassword == "" && req.NewPassword == "" {
		return errors.New("Not enter the New Or old password")
	}
	
	if err := utils.Comparepassword(user.Password, req.NewPassword); err != nil {
		return errors.New("Old Password Not Matching")
	}

	newHashed, err := utils.Hashing(req.NewPassword)
	if err != nil {
		return err
	} 

	field := make(map[string]interface{})
	field["password"] = newHashed
	if err := s.repo.Update(&User{}, field, "id = ?", userID); err != nil {
		return err
	}

	return nil
}

// Delete Account
func (s *ProfileService) DeleteAccount(userID uint) error {
	return s.repo.Delete(&User{}, userID)
}