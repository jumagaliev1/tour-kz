package postgre

import (
	"context"
	"gorm.io/gorm"
	"tour-kz/internal/logger"
	"tour-kz/internal/model"
)

type AccountRepository struct {
	DB     *gorm.DB
	logger logger.RequestLogger
}

func NewAccountRepository(DB *gorm.DB, logger logger.RequestLogger) *AccountRepository {
	return &AccountRepository{DB: DB, logger: logger}
}

func (r *AccountRepository) Create(ctx context.Context, account model.Account) (uint, error) {
	if err := r.DB.WithContext(ctx).Create(&account).Error; err != nil {
		return 0, err
	}

	return account.ID, nil
}

func (r *AccountRepository) GetByUserID(ctx context.Context, userID uint) (*model.Account, error) {
	var account *model.Account
	if err := r.DB.WithContext(ctx).Find(&account, "user_id = ?", userID).Scan(&account).Error; err != nil {
		return nil, err
	}

	return account, nil
}

func (r *AccountRepository) Update(ctx context.Context, account model.Account) error {
	if err := r.DB.WithContext(ctx).Save(&account).Error; err != nil {
		return err
	}

	return nil
}

func (r *AccountRepository) UpdateFirstLevel(ctx context.Context, userID uint, amount int) error {
	err := r.DB.WithContext(ctx).Raw(
		`UPDATE accounts a
				SET balance = balance + (? * 0.2)
				WHERE a.user_id in (SELECT u.id from users u JOIN referrals r on r.parent_id = u.id and r.user_id = ?)`, amount, userID).Scan(&model.Account{}).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *AccountRepository) UpdateSecondLevel(ctx context.Context, userID uint, amount int) error {
	err := r.DB.WithContext(ctx).Raw(
		`UPDATE accounts a
				SET balance = a.balance + (? * 0.02)
				WHERE a.user_id in (SELECT u.id from users u JOIN referrals r on r.parent_id = u.id and r.user_id
				 in (SELECT u.id from users u JOIN referrals r on r.parent_id = u.id and r.user_id = ?))`, amount, userID).Scan(&model.Account{}).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *AccountRepository) UpdateThirdLevel(ctx context.Context, userID uint, amount int) error {
	err := r.DB.WithContext(ctx).Raw(
		`UPDATE accounts a
			SET balance = a.balance + (? * 0.01)
			WHERE a.user_id in (SELECT u.id from users u join referrals r on r.parent_id = u.id and r.user_id in (SELECT u.id from users u JOIN referrals r on r.parent_id = u.id and r.user_id in (SELECT u.id from users u JOIN referrals r on r.parent_id = u.id and r.user_id = ?)))`, amount, userID).Scan(&model.Account{}).Error

	if err != nil {
		return err
	}

	return nil
}
