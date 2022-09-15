package guest

import (
	"net/http"

	"github.com/khotchapan/KonLakRod-api/internal/core/context"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service GuestInterface
}

func NewHandler(service GuestInterface) *Handler {
	return &Handler{
		service: service,
		//service: NewService(c),
	}
}
func (h *Handler) LoginUsers(c echo.Context) error {
	// username := c.FormValue("username")
	// password := c.FormValue("password")
	request := &LoginUsersForm{}
	cc := c.(*context.CustomContext)
	if err := cc.BindAndValidate(request); err != nil {
		return err
	}
	t, err := h.service.LoginUsers(c, request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}
