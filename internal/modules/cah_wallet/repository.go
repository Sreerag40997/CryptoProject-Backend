package cashwallet

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
  Create(ctx context.Context, wallet *Wallet) error
	GetByUserID(ctx context.Context, userID uint) (*Wallet, error)
}
type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{
		db: db,
	}
}


func (r *repo) Create(ctx context.Context, wallet *Wallet) error {
	return r.db.WithContext(ctx).Create(wallet).Error
}

func (r *repo) GetByUserID(ctx context.Context, userID uint) (*Wallet, error) {
	var wallet Wallet

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&wallet).Error

	if err != nil {
		return nil, err
	}

	return &wallet, nil
}