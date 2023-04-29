package middleware

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"tour-kz/internal/service"
)

const (
	ADMIN_ROLE  = "admin"
	CLIENT_ROLE = "client"
)

var ErrAdminRequire = errors.New("you should be admin")

type Permission struct {
	User service.IUserService
}

func NewPermission(user service.IUserService) *Permission {
	return &Permission{User: user}
}

func (m *Permission) AdminRequire(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := m.User.GetUserFromRequest(c.Request().Context())
		if err != nil {
			return err
		}

		if user.Role != ADMIN_ROLE {
			return c.JSON(http.StatusForbidden, ErrAdminRequire.Error())
		}

		return next(c)
	}
}
