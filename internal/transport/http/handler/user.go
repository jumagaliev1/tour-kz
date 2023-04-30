package handler

import (
	"github.com/labstack/echo/v4"
	_ "github.com/swaggo/echo-swagger"
	"net/http"
	"tour-kz/internal/logger"
	"tour-kz/internal/model"
	"tour-kz/internal/service"
	jwt "tour-kz/internal/transport/middleware"
)

type UserHandler struct {
	service *service.Manager
	jwt     *jwt.JWTAuth
	logger  logger.RequestLogger
}

func NewUserHandler(service *service.Manager, jwt *jwt.JWTAuth, logger logger.RequestLogger) *UserHandler {
	return &UserHandler{
		service: service,
		jwt:     jwt,
		logger:  logger,
	}
}

// CreateUser godoc
// @Summary      Создание пользователя
// @Description  Создание пользователя
// @ID           CreateUser
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        rq   body      model.UserCreateReq  true  "Input body"
// @Success	     200  {object}  model.User
// @Router       /user [post]
func (h *UserHandler) Create(c echo.Context) error {
	h.logger.Logger(c.Request().Context()).Info("creating user...")
	var user model.UserCreateReq
	if err := c.Bind(&user); err != nil {
		h.logger.Logger(c.Request().Context()).Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	usr, err := h.service.User.Create(c.Request().Context(), user)
	if err != nil {
		switch err {
		case model.ErrDuplicateEmail:
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		h.logger.Logger(c.Request().Context()).Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, usr)
}

// Auth godoc
// @Summary      Auth get JWT token
// @Description Auth get JWT token
// @ID           AuthUser
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        rq   body      model.AuthUser  true  "Входящие данные"
// @Success	     200  {object}  string
// @Router       /tokens/authentication [post]
func (h *UserHandler) Auth(c echo.Context) error {
	var input model.AuthUser

	if err := c.Bind(&input); err != nil {
		h.logger.Logger(c.Request().Context()).Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	h.logger.Logger(c.Request().Context()).Info(input.Password)
	if err := h.service.User.Auth(c.Request().Context(), input); err != nil {
		h.logger.Logger(c.Request().Context()).Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	token, err := h.jwt.GenerateJWT(input.Email)
	if err != nil {
		h.logger.Logger(c.Request().Context()).Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, token)
}

// Get User godoc
// @Summary      Get User
// @Description  Get User
// @Security	ApiKeyAuth
// @ID           GetUser
// @Tags         user
// @Accept       json
// @Produce      json
// @Success	     200  {object}  model.User
// @Router       /user [get]
func (h *UserHandler) Get(c echo.Context) error {
	user, err := h.service.User.GetUserFromRequest(c.Request().Context())
	if err != nil {
		h.logger.Logger(c.Request().Context()).Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

// GetAll User godoc
// @Summary      Get All User
// @Description  Get All User
// @Security	ApiKeyAuth
// @ID           GetAllUser
// @Tags         user
// @Accept       json
// @Produce      json
// @Success	     200  {object}  []model.User
// @Router       /users [get]
func (h *UserHandler) GetAll(c echo.Context) error {
	users, err := h.service.User.GetAll(c.Request().Context())
	if err != nil {
		h.logger.Logger(c.Request().Context()).Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, users)
}
