package ecard



import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, card *Card) error
	GetByUserID(ctx context.Context, userID uint) (*Card, error)
	UpdateStatus(ctx context.Context, userID uint, status,blockedBy string) error 
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, card *Card) error {
	return r.db.WithContext(ctx).Create(card).Error
}

func (r *repo) GetByUserID(ctx context.Context, userID uint) (*Card, error) {
	var card Card
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&card).Error

	if err != nil {
		return nil, err
	}

	return &card, nil
}

func (r *repo) UpdateStatus(ctx context.Context, userID uint, status, blockedBy string) error {
	return r.db.WithContext(ctx).
		Model(&Card{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"status":status,
			"blocked_by":blockedBy,
		}).Error
}