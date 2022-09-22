package guest

import (
	"net/http"

	"github.com/khotchapan/KonLakRod-api/internal/middleware"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service ServiceInterface
}

func NewHandler(service ServiceInterface) *Handler {
	return &Handler{
		service: service,
	}
}
func (h *Handler) LoginUsers(c echo.Context) error {
	request := &LoginUsersForm{}
	cc := c.(*middleware.CustomContext)
	if err := cc.BindAndValidate(request); err != nil {
		return err
	}
	t, err := h.service.LoginUsers(c, request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, t)
}
