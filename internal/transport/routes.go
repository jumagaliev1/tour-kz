package transport

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
)

func (s *Server) SetupRoutes() *echo.Group {
	v1 := s.App.Group("/api/v1")
	s.App.GET("/ready", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	s.App.GET("/live", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	v1.POST("/user", s.handler.User.Create)
	v1.POST("/tokens/authentication", s.handler.User.Auth)
	v1.GET("/user", s.handler.User.Get, s.jwt.ValidateAuth)
	v1.GET("/users", s.handler.User.GetAll)

	v1.GET("/my_referrals", s.handler.Referral.GetReferrals, s.jwt.ValidateAuth)
	v1.POST("/add_balance", s.handler.Account.AddBalance, s.jwt.ValidateAuth, s.permission.AdminRequire)
	v1.GET("/my_balance", s.handler.Account.MyBalance, s.jwt.ValidateAuth)

	v1.POST("/payment/income", s.handler.Payment.CreateIncome, s.jwt.ValidateAuth)
	v1.POST("/payment/outcome", s.handler.Payment.CreateOutcome, s.jwt.ValidateAuth)
	v1.GET("/payments", s.handler.Payment.GetPayments, s.jwt.ValidateAuth, s.permission.AdminRequire)

	v1.PATCH("/payment/:id", s.handler.Payment.ApprovePayment, s.jwt.ValidateAuth, s.permission.AdminRequire)
	v1.DELETE("/payment/:id", s.handler.Payment.CancelPayment, s.jwt.ValidateAuth, s.permission.AdminRequire)

	s.App.GET("/swagger/*", echoSwagger.WrapHandler)

	return v1
}
