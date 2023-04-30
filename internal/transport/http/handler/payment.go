package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"tour-kz/internal/logger"
	"tour-kz/internal/model"
	"tour-kz/internal/service"
)

type PaymentHandler struct {
	service *service.Manager
	logger  logger.RequestLogger
}

func NewPaymentHandler(service *service.Manager, logger logger.RequestLogger) *PaymentHandler {
	return &PaymentHandler{service: service, logger: logger}
}

// CreateIncome godoc
// @Summary      CreateIncome Payment
// @Description  CreateIncome Payment for checking
// @ID           CreateIncomePayment
// @Tags         payment
// @Security	ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        rq   body      model.TypeCreateReq  true  "Input body"
// @Success	     201  {object}  int
// @Router       /payment/income [post]
func (h *PaymentHandler) CreateIncome(c echo.Context) error {
	var body model.TypeCreateReq
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user, err := h.service.User.GetUserFromRequest(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}

	body.UserID = user.ID
	body.Type = model.TypeIncome

	id, err := h.service.Payment.Create(c.Request().Context(), body.ToPayment())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, id)
}

// CreateOutcome godoc
// @Summary      CreateOutcome Payment
// @Description  CreateOutcome Payment for checking
// @ID           CreateOutcomePayment
// @Tags         payment
// @Security	ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        rq   body      model.TypeCreateReq  true  "Input body"
// @Success	     201  {object}  int
// @Router       /payment/outcome [post]
func (h *PaymentHandler) CreateOutcome(c echo.Context) error {
	var body model.TypeCreateReq
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user, err := h.service.User.GetUserFromRequest(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}

	body.UserID = user.ID
	body.Type = model.TypeOutcome

	id, err := h.service.Payment.Create(c.Request().Context(), body.ToPayment())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, id)
}

// GetPayments  godoc
// @Summary      Get Payments
// @Description  Get Payments with sort
// @Security	ApiKeyAuth
// @ID           GetPayments
// @Tags         payment
// @Accept       json
// @Produce      json
// @Success	     200  {object}  []model.Payment
// @Router       /payments [get]
func (h *PaymentHandler) GetPayments(c echo.Context) error {
	resp, err := h.service.Payment.GetPayments(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, resp)
}

// ApprovePayment  godoc
// @Summary      Approve Payment
// @Description  GApprovePayment
// @Security	ApiKeyAuth
// @ID           ApprovePayment
// @Tags         payment
// @Accept       json
// @Produce      json
// @Param        id   path      int  true   "input"
// @Success	     200  {object}  []model.Payment
// @Router       /payment/{id} [patch]
func (h *PaymentHandler) ApprovePayment(c echo.Context) error {
	var ID uint
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Logger(c.Request().Context()).Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}
	ID = uint(id)
	payment, err := h.service.Payment.GetByID(c.Request().Context(), ID)
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, err)
	}

	account, err := h.service.Account.GetByUser(c.Request().Context(), payment.UserID)
	if err != nil {
		return err
	}

	h.logger.Logger(c.Request().Context()).Info("payment:", payment.Amount, "balance:", account.Balance)

	if payment.Type == model.TypeIncome {
		err = h.service.Account.UpdateLevels(c.Request().Context(), payment.UserID, payment.Amount)
		if err != nil {
			return err
		}
	} else if payment.Type == model.TypeOutcome {
		if account.Balance-payment.Amount < 50000 {
			return c.JSON(http.StatusServiceUnavailable, "not enough money")
		}
		account.Balance = account.Balance - payment.Amount

		err = h.service.Account.Update(c.Request().Context(), *account)
		if err != nil {
			return err
		}
	}

	h.logger.Logger(c.Request().Context()).Info("After:", account.Balance)

	payment.Status = model.StatusComplete

	err = h.service.Payment.Update(c.Request().Context(), *payment)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, payment)
}

// CancelPayment  godoc
// @Summary      Cancel Payment
// @Description  Cancel Payment
// @Security	ApiKeyAuth
// @ID           CancelPayment
// @Tags         payment
// @Accept       json
// @Produce      json
// @Param        id   path      int  true   "input"
// @Success	     200  {object}  []model.Payment
// @Router       /payment/{id} [delete]
func (h *PaymentHandler) CancelPayment(c echo.Context) error {
	var ID uint
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Logger(c.Request().Context()).Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}
	ID = uint(id)
	payment, err := h.service.Payment.GetByID(c.Request().Context(), ID)
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, err)
	}

	payment.Status = model.StatusCancel

	err = h.service.Payment.Update(c.Request().Context(), *payment)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, payment)
}
