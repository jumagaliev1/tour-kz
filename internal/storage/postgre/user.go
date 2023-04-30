package postgre

import (
	"context"
	"database/sql"
	"errors"
	"gorm.io/gorm"
	"tour-kz/internal/logger"
	"tour-kz/internal/model"
)

type UserRepository struct {
	DB     *gorm.DB
	logger logger.RequestLogger
}

func NewUserRepository(DB *gorm.DB, logger logger.RequestLogger) *UserRepository {
	return &UserRepository{DB: DB, logger: logger}
}

func (r *UserRepository) Create(ctx context.Context, user model.User) (*model.User, error) {
	tx := r.DB.Begin()

	if err := r.DB.WithContext(ctx).Create(&user).Error; err != nil {
		r.logger.Logger(ctx).Error(err)
		tx.Rollback()
		switch {
		case model.IsDuplicateError(err):
			return nil, model.ErrDuplicateEmail
		default:
			return nil, err
		}
	}

	account := &model.Account{
		UserID:  user.ID,
		User:    user,
		Balance: 0,
	}

	if err := r.DB.WithContext(ctx).Create(&account).Error; err != nil {
		r.logger.Logger(ctx).Error(err)
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user model.User) error {
	if err := r.DB.WithContext(ctx).Save(user).Error; err != nil {
		r.logger.Logger(ctx).Error(err)
		return err
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, ID int) error {
	if err := r.DB.WithContext(ctx).Delete(model.User{}, ID).Error; err != nil {
		r.logger.Logger(ctx).Error(err)
		return err
	}

	return nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]*model.User, error) {
	var users []*model.User

	if err := r.DB.WithContext(ctx).Find(&users).Error; err != nil {
		r.logger.Logger(ctx).Error(err)
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, username string) (*model.User, error) {
	var user model.User

	if err := r.DB.WithContext(ctx).Where("email = ?", username).First(&user).Error; err != nil {
		r.logger.Logger(ctx).Error(err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, model.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, ID uint) (*model.User, error) {
	var user model.User

	if err := r.DB.WithContext(ctx).Where("id = ?", ID).First(&user).Error; err != nil {
		r.logger.Logger(ctx).Error(err)
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByReferralCode(ctx context.Context, referralCode string) (*model.User, error) {
	var user model.User

	if err := r.DB.WithContext(ctx).Where("referral_code = ?", referralCode).First(&user).Error; err != nil {
		r.logger.Logger(ctx).Error(err)
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, errors.New("invalid referral code")
		default:
			return nil, err
		}
	}

	return &user, nil
}
