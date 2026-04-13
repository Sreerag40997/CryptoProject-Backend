package profile

import "time"

type User struct {
	ID            uint   `gorm:"primaryKey"`
	Name          string `gorm:"not null"`
	Email         string `gorm:"unique;not null"`
	Password      string `gorm:"not null"`
	Role          string `gorm:"default:user"`
	KYCStatus     bool
	IsVerified    bool
	IsBlocked     bool
	ProfilePicURL string    `json:"profile_pic_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}