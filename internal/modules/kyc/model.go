package kyc

import "time"

type KYC struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint `gorm:"uniqueIndex;not null"`

	FullName string
	DOB      time.Time

	AccountNumber string `gorm:"type:text"`
	AccountLast4  string
	IFSC   string

	AadhaarNumber string `gorm:"type:text"`
	PANNumber     string `gorm:"type:text"`

	AadhaarFrontURL string
	AadhaarBackURL  string
	PANURL          string
	SelfieURL       string

	Status string `gorm:"type:varchar(20);default:'pending';index"`

	RejectionReason string
	VerifiedAt      *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}