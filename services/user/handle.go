package user

import (
	"net/http"

	"github.com/khotchapan/KonLakRod-api/internal/core/context"
	"github.com/khotchapan/KonLakRod-api/mongodb"
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
		//return echo.NewHTTPError(http.StatusBadRequest, err)
		return c.JSON(http.StatusBadRequest, err)
	}
	response, err := h.service.FindAllUsers(c, request)
	if err != nil {
		//return echo.NewHTTPError(http.StatusBadRequest, err)
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, response)
}

func (h *Handler) GetOneUsers(c echo.Context) error {
	request := &GetOneUsersForm{}
	cc := c.(*context.CustomContext)
	if err := cc.BindAndValidate(request); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	response, err := h.service.FindOneUsers(c, request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, response)
}

func (h *Handler) PostUsers(c echo.Context) error {
	request := &CreateUsersForm{}
	cc := c.(*context.CustomContext)
	if err := cc.BindAndValidate(request); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	err := h.service.CreateUsers(c, request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	response := &mongodb.Response{}
	return c.JSON(http.StatusOK, response.SuccessfulCreated())
}
func (h *Handler) PutUsers(c echo.Context) error {
	request := &UpdateUsersForm{}
	cc := c.(*context.CustomContext)
	if err := cc.BindAndValidate(request); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	err := h.service.UpdateUsers(c, request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	response := &mongodb.Response{}
	return c.JSON(http.StatusOK, response.SuccessfulOK())
}

func (h *Handler) DeleteUsers(c echo.Context) error {
	request := &DeleteUsersForm{}
	cc := c.(*context.CustomContext)
	if err := cc.BindAndValidate(request); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	err := h.service.DeleteUsers(c, request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	response := &mongodb.Response{}
	return c.JSON(http.StatusOK, response.SuccessfulOK())
}

func (h *Handler) UploadFile(c echo.Context) error {
	var req UploadForm
	file, _ := c.FormFile("file")
	req.File = file
	res, err := h.service.UploadFile(c, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"link": res,
	})
	
}
