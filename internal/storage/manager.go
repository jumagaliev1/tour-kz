package storage

import (
	"context"
	"tour-kz/internal/config"
	"tour-kz/internal/logger"
	"tour-kz/internal/storage/postgre"
)

type Manager struct {
	User     IUserRepository
	Referral IReferralRepository
	Account  IAccountRepository
	Payment  IPaymentRepository
}

func New(ctx context.Context, cfg *config.Config, logger logger.RequestLogger) (*Manager, error) {
	pgDB, err := postgre.Dial(ctx, cfg.Postgres)
	if err != nil {
		return nil, err
	}
	userRepo := postgre.NewUserRepository(pgDB, logger)
	referRepo := postgre.NewReferralRepository(pgDB, logger)
	accountRepo := postgre.NewAccountRepository(pgDB, logger)
	paymentRepo := postgre.NewPaymentRepository(pgDB, logger)
	return &Manager{
		User:     userRepo,
		Referral: referRepo,
		Account:  accountRepo,
		Payment:  paymentRepo,
	}, nil
}
