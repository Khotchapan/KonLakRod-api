package user

import (
	"net/http"

	"github.com/khotchapan/KonLakRod-api/internal/core/context"
	"github.com/khotchapan/KonLakRod-api/mongodb/user"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service UserInterface
}

func NewHandler(service UserInterface) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetAllUsers(c echo.Context) error {
	request := &user.GetAllUsersForm{}
	cc := c.(*context.CustomContext)
	if err := cc.BindAndValidate(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	response, err := h.service.FindAllUsers(c, request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, response)
}
