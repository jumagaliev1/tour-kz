package postgre

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"tour-kz/internal/logger"
	"tour-kz/internal/model"
)

type ReferralRepository struct {
	DB     *gorm.DB
	logger logger.RequestLogger
}

func NewReferralRepository(DB *gorm.DB, logger logger.RequestLogger) *ReferralRepository {
	return &ReferralRepository{DB: DB, logger: logger}
}

func (r *ReferralRepository) Create(ctx context.Context, referral model.Referral) (uint, error) {
	if err := r.DB.WithContext(ctx).Create(&referral).Error; err != nil {
		r.logger.Logger(ctx).Error(err)
		switch {
		case errors.Is(err, gorm.ErrDuplicatedKey):
			return 0, model.ErrDuplicateKey
		default:
			return 0, err
		}
	}

	return referral.ID, nil
}

func (r *ReferralRepository) GetFirstLine(ctx context.Context, userID uint) ([]model.User, error) {
	var users []model.User
	err := r.DB.WithContext(ctx).Raw("SELECT u.id, u.first_name, u.last_name, u.email, u.password, u.phone, u.referral_code FROM users u JOIN referrals r on r.user_id = u.id and r.parent_id = ?", userID).Scan(&users).Error
	if err != nil {
		r.logger.Logger(ctx).Error(err)
		return nil, err
	}

	return users, nil
}

func (r *ReferralRepository) GetSecondLine(ctx context.Context, userID uint) ([]model.User, error) {
	var users []model.User
	err := r.DB.WithContext(ctx).Raw(
		`SELECT u.id, u.first_name, u.last_name, u.email, u.password, u.phone, u.referral_code 
				FROM users u JOIN referrals r on r.user_id = u.id and r.parent_id in 
				            (SELECT u.id FROM users u JOIN referrals r on r.user_id = u.id and r.parent_id = ?)`, userID).Scan(&users).Error
	if err != nil {
		r.logger.Logger(ctx).Error(err)
		return nil, err
	}

	return users, nil
}

func (r *ReferralRepository) GetThirdLine(ctx context.Context, userID uint) ([]model.User, error) {
	var users []model.User
	err := r.DB.WithContext(ctx).Raw(
		`SELECT u.id, u.first_name, u.last_name, u.email, u.password, u.phone, u.referral_code 
			from users u JOIN referrals r on r.user_id = u.id and r.parent_id in 
			    (SELECT u.id FROM users u JOIN referrals r on r.user_id = u.id and r.parent_id in 
			    (SELECT u.id FROM users u JOIN referrals r on r.user_id = u.id and r.parent_id = ?))`, userID).Scan(&users).Error
	if err != nil {
		r.logger.Logger(ctx).Error(err)
		return nil, err
	}

	return users, nil
}
