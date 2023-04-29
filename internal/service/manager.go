package service

import (
	"errors"
	"tour-kz/internal/config"
	"tour-kz/internal/logger"
	"tour-kz/internal/storage"
)

type Manager struct {
	User     IUserService
	Referral IReferralService
	Account  IAccountService
	Payment  IPaymentService
}

func New(repo *storage.Manager, cfg *config.Config, logger logger.RequestLogger) (*Manager, error) {
	if repo == nil {
		return nil, errors.New("No storage")
	}
	userService := NewUserService(repo, cfg, logger)
	referService := NewReferralService(repo, logger)
	accountService := NewAccountService(repo, logger)
	paymentService := NewPaymentService(repo, logger)
	return &Manager{
		User:     userService,
		Referral: referService,
		Account:  accountService,
		Payment:  paymentService,
	}, nil
}
