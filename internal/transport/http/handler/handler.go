package handler

import "github.com/labstack/echo/v4"

type IUserHandler interface {
	Create(c echo.Context) error
	Get(c echo.Context) error
	Auth(c echo.Context) error
	GetAll(c echo.Context) error
}

type IReferralHandler interface {
	GetReferrals(c echo.Context) error
}

type IAccountHandler interface {
	AddBalance(c echo.Context) error
	MyBalance(c echo.Context) error
}

type IPaymentHandler interface {
	CreateIncome(c echo.Context) error
	CreateOutcome(c echo.Context) error
	GetPayments(c echo.Context) error
	ApprovePayment(c echo.Context) error
	CancelPayment(c echo.Context) error
}
