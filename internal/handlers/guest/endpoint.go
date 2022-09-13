package guest

import (
	"log"
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

func (h *Handler) PostLoginUsers(c echo.Context) error {
	log.Println("========STEP1========")
	request := &LoginUsersForm{}
	cc := c.(*context.CustomContext)
	if err := cc.BindAndValidate(request); err != nil {
		return err
	}
	response, err := h.service.LoginUsers(c, request)
	if err != nil {
		log.Println("endpoint")
		log.Println("err", err)
		return  echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	//response := &mongodb.Response{}
	//return c.JSON(http.StatusOK, response.SuccessfulCreated())
	return c.JSON(http.StatusOK, response)
}
