package service

import (
	"context"
	"tour-kz/internal/model"
)

type IUserService interface {
	Create(ctx context.Context, user model.UserCreateReq) (*model.User, error)
	Update(ctx context.Context, user model.User) error
	Delete(ctx context.Context, ID int) error
	GetAll(ctx context.Context) ([]*model.User, error)
	CheckPassword(encPass, providedPassword string) error
	HashPassword(password string) (string, error)
	Auth(ctx context.Context, user model.AuthUser) error
	RefreshToken() (string, error)
	GenerateToken(user model.AuthUser) (string, error)
	ParseToken(accessToken string) (string, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserFromRequest(ctx context.Context) (*model.User, error)
}

type IReferralService interface {
	Create(ctx context.Context, referral model.Referral) (uint, error)
	GetFirstLine(ctx context.Context, userID uint) ([]model.User, error)
	GetSecondLine(ctx context.Context, userID uint) ([]model.User, error)
	GetThirdLine(ctx context.Context, userID uint) ([]model.User, error)
	GetReferrals(ctx context.Context, userID uint) (map[string][]model.User, error)
}

type IAccountService interface {
	Create(ctx context.Context, account model.Account) (uint, error)
	GetByUser(ctx context.Context, userID uint) (*model.Account, error)
	Update(ctx context.Context, account model.Account) error
	UpdateLevels(ctx context.Context, userID uint, amount int) error
}

type IPaymentService interface {
	Create(ctx context.Context, payment model.Payment) (uint, error)
	GetByID(ctx context.Context, ID uint) (*model.Payment, error)
	Update(ctx context.Context, payment model.Payment) error
	GetPayments(ctx context.Context) (map[string][]model.Payment, error)
}
