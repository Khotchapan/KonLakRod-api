package test

import (
	"net/http"

	"github.com/khotchapan/KonLakRod-api/internal/core/context"
	googleCloud "github.com/khotchapan/KonLakRod-api/lagacy/google/google_cloud"
	"github.com/khotchapan/KonLakRod-api/mongodb"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service TestInterface
}

func NewHandler(service TestInterface) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetFile(c echo.Context) error {
	// var req UploadForm
	// file, _ := c.FormFile("file")
	// req.File = file
	response, err := h.service.FindAllBooks(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// return c.JSON(http.StatusOK, map[string]interface{}{
	// 	"link": res,
	// })
	return c.JSON(http.StatusOK, response)
}

func (h *Handler) GetOneGoogleCloudBooks(c echo.Context) error {
	request := &GetOneGoogleCloudBooksForm{}
	cc := c.(*context.CustomContext)
	if err := cc.BindAndValidate(request); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	response, err := h.service.FindOneBooks(c, request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, response)
}
func (h *Handler) PostGoogleCloudBooks(c echo.Context) error {
	request := &googleCloud.CreateBooksForm{}
	cc := c.(*context.CustomContext)

	if err := cc.BindAndValidate(request); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	err := h.service.CreateBooks(c, request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	response := &mongodb.Response{}
	return c.JSON(http.StatusOK, response.SuccessfulCreated())
}
