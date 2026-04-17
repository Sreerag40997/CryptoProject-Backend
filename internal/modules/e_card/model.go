package ecard



import "time"

type Card struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint `gorm:"uniqueIndex"`

	CardNumber string 
	Last4      string

	ExpiryMonth int
	ExpiryYear  int

	CVV string 

	Status string `gorm:"default:'active'"`
	BlockedBy string

	CreatedAt time.Time
	UpdatedAt time.Time
}