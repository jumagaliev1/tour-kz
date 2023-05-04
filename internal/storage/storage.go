package storage

import (
	"context"
	"tour-kz/internal/model"
)

type IUserRepository interface {
	Create(ctx context.Context, user model.User) (*model.User, error)
	Update(ctx context.Context, user model.User) error
	Delete(ctx context.Context, ID int) error
	GetAll(ctx context.Context) ([]*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, ID uint) (*model.User, error)
	GetByReferralCode(ctx context.Context, referralCode string) (*model.User, error)
}

type IReferralRepository interface {
	Create(ctx context.Context, referral model.Referral) (uint, error)
	GetFirstLine(ctx context.Context, userID uint) ([]model.User, error)
	GetSecondLine(ctx context.Context, userID uint) ([]model.User, error)
	GetThirdLine(ctx context.Context, userID uint) ([]model.User, error)
}

type IAccountRepository interface {
	Create(ctx context.Context, account model.Account) (uint, error)
	Update(ctx context.Context, account model.Account) error
	GetByUserID(ctx context.Context, userID uint) (*model.Account, error)
	UpdateBalance(ctx context.Context, userID uint, amount int) error
	UpdateFirstLevel(ctx context.Context, userID uint, amount int) error
	UpdateSecondLevel(ctx context.Context, userID uint, amount int) error
	UpdateThirdLevel(ctx context.Context, userID uint, amount int) error
}

type IPaymentRepository interface {
	Create(ctx context.Context, payment model.Payment) (uint, error)
	Update(ctx context.Context, payment model.Payment) error
	GetByID(ctx context.Context, ID uint) (*model.Payment, error)
	GetByNonCompleted(ctx context.Context) ([]model.Payment, error)
	GetByIncome(ctx context.Context) ([]model.Payment, error)
	GetByOutcome(ctx context.Context) ([]model.Payment, error)
	GetAll(ctx context.Context) ([]model.Payment, error)
}
