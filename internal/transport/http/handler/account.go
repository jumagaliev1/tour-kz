package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"tour-kz/internal/logger"
	"tour-kz/internal/model"
	"tour-kz/internal/service"
)

type AccountHandler struct {
	service *service.Manager
	logger  logger.RequestLogger
}

func NewAccountHandler(service *service.Manager, logger logger.RequestLogger) *AccountHandler {
	return &AccountHandler{service: service, logger: logger}
}

// AddBalance godoc
// @Summary      Add balance for User
// @Description  Add balance for User
// @Security	ApiKeyAuth
// @ID           AddBalance
// @Tags         account
// @Accept       json
// @Produce      json
// @Param        rq   body      model.AddBalanceReq  true  "Input body"
// @Success	     200  {object}  string
// @Router       /add_balance [post]
func (h *AccountHandler) AddBalance(c echo.Context) error {
	var body model.AddBalanceReq

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	h.logger.Logger(c.Request().Context()).Info(body)

	err := h.service.Account.UpdateLevels(c.Request().Context(), body.UserID, body.Amount)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "Succesfully")
}

// MyBalance  godoc
// @Summary      Get user balance
// @Description  Get authorized user balance
// @Security	ApiKeyAuth
// @ID           MyBalance
// @Tags         account
// @Accept       json
// @Produce      json
// @Success	     200  {object}  int
// @Router       /my_balance [get]
func (h *AccountHandler) MyBalance(c echo.Context) error {
	user, err := h.service.User.GetUserFromRequest(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	account, err := h.service.Account.GetByUser(c.Request().Context(), user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]int{"my_balance": account.Balance})
}
