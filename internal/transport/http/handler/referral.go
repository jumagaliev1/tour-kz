package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"tour-kz/internal/logger"
	_ "tour-kz/internal/model"
	"tour-kz/internal/service"
)

type ReferralHandler struct {
	service *service.Manager
	logger  logger.RequestLogger
}

func NewReferralHandler(service *service.Manager, logger logger.RequestLogger) *ReferralHandler {
	return &ReferralHandler{
		service: service,
		logger:  logger,
	}
}

// Get Referrals godoc
// @Summary      Get Referrals
// @Description  Get Referrals
// @Security	ApiKeyAuth
// @ID           GetUReferral
// @Tags         referral
// @Accept       json
// @Produce      json
// @Success	     200  {object}  model.User
// @Router       /my_referrals [get]
func (h *ReferralHandler) GetReferrals(c echo.Context) error {
	user, err := h.service.User.GetUserFromRequest(c.Request().Context())
	if err != nil {
		h.logger.Logger(c.Request().Context()).Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	users, err := h.service.Referral.GetReferrals(c.Request().Context(), user.ID)
	if err != nil {
		h.logger.Logger(c.Request().Context()).Error(err)
		return echo.NewHTTPError(http.StatusNoContent, err.Error())
	}

	return c.JSON(http.StatusOK, users)
}
