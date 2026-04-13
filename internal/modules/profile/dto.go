package profile

import "time"

type UserProfile struct {
	ID            uint   
	Name          string 
	Email         string 
	Role          string 
	KYCStatus     bool
	IsVerified    bool
	IsBlocked     bool
	ProfilePicURL string    `json:"profile_pic_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type EditProfileRequest struct {
	Name string `json:"newname"`
}

type ChangePassWordReq struct {
	OldPassword string `json:"oldpassword" validate:"required"`
	NewPassword string `json:"newpassword" validate:"required"`
}
