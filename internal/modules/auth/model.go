package auth

import "time"

type User struct {
	ID            uint   `gorm:"primaryKey"`
	Name          string `gorm:"not null"`
	Email         string `gorm:"unique;not null"`
	Password      string `gorm:"not null"`
	Role          string `gorm:"default:user"`
	KYCStatus     bool	`gorm:"default:false"`
	IsVerified    bool	`gorm:"default:false"`
	IsBlocked     bool	`gorm:"default:false"`
	ProfilePicURL string    `json:"profile_pic_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

//KYCStatus     string `gorm:"type:varchar(20);default:'pending';index"`
