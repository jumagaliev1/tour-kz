package handler

import (
	"errors"
	"tour-kz/internal/logger"
	"tour-kz/internal/service"
	jwt "tour-kz/internal/transport/middleware"
)

type Manager struct {
	User     IUserHandler
	Referral IReferralHandler
	Account  IAccountHandler
	Payment  IPaymentHandler
}

func New(service *service.Manager, jwt *jwt.JWTAuth, logger logger.RequestLogger) (*Manager, error) {
	if service == nil {
		return nil, errors.New("No given service")
	}

	user := NewUserHandler(service, jwt, logger)
	referral := NewReferralHandler(service, logger)
	account := NewAccountHandler(service, logger)
	payment := NewPaymentHandler(service, logger)
	return &Manager{
		User:     user,
		Referral: referral,
		Account:  account,
		Payment:  payment,
	}, nil
}
