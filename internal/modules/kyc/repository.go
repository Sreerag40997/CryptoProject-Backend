package kyc

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	// User Side
	Create(ctx context.Context, kyc *KYC) error
	GetByUserID(ctx context.Context, userID uint) (*KYC, error)
	Update(ctx context.Context, kyc *KYC) error

	// Admin Side
	GetByID(ctx context.Context, id uint) (*KYC, error)
	ListPending(ctx context.Context,status string) ([]KYC, error)
	UpdateStatus(ctx context.Context, id uint, status string, reason string) error
	ListWithFilter(ctx context.Context, status string, limit, offset int) ([]KYC, error)
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{
		db: db,
	}
}

func (r *repo) Create(ctx context.Context, kyc *KYC) error {
	return r.db.WithContext(ctx).Create(kyc).Error
}

func (r *repo) GetByUserID(ctx context.Context, userID uint) (*KYC, error) {
	var kyc KYC
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&kyc).Error

	if err != nil {
		return nil, err
	}

	return &kyc, nil
}

func (r *repo) Update(ctx context.Context, kyc *KYC) error {
	return r.db.WithContext(ctx).Save(kyc).Error
}

func (r *repo) GetByID(ctx context.Context, id uint) (*KYC, error) {
	var kyc KYC

	err := r.db.WithContext(ctx).
		First(&kyc, id).Error

	if err != nil {
		return nil, err
	}

	return &kyc, nil
}

func (r *repo) ListPending(ctx context.Context,status string) ([]KYC, error) {
	var kycs []KYC

	query := r.db.WithContext(ctx)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.
		Order("created_at DESC").
		Find(&kycs).Error

	return kycs, err
}

func (r *repo) UpdateStatus(ctx context.Context, id uint, status string, reason string) error {
	updateData := map[string]interface{}{
		"status": status,
	}

	if status == "rejected" {
		updateData["rejection_reason"] = reason
	}

	if status == "approved" {
		updateData["rejection_reason"] = ""
		updateData["verified_at"] = gorm.Expr("NOW()")
	}

	return r.db.WithContext(ctx).
		Model(&KYC{}).
		Where("id = ?", id).
		Updates(updateData).Error
}

func (r *repo) ListWithFilter(ctx context.Context, status string, limit, offset int) ([]KYC, error) {
	var kycs []KYC

	query := r.db.WithContext(ctx)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&kycs).Error

	return kycs, err
}