package auth

import "time"

type UserRegisterRequest struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type UserRegisterResponse struct {
	Name string `json:"name"`
	Email string `json:"email"`
}

type UserLoginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
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

