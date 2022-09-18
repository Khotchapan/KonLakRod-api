package test

import (
	"net/http"

	"github.com/khotchapan/KonLakRod-api/internal/middleware"
	"github.com/khotchapan/KonLakRod-api/internal/core/mongodb"
	googleCloud "github.com/khotchapan/KonLakRod-api/internal/lagacy/google/google_cloud"
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
	cc := c.(*middleware.CustomContext)
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
	cc := c.(*middleware.CustomContext)

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

func (h *Handler) PutBooks(c echo.Context) error {
	request := &googleCloud.UpdateBooksForm{}
	cc := c.(*middleware.CustomContext)
	if err := cc.BindAndValidate(request); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	err := h.service.UpdateBooks(c, request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	response := &mongodb.Response{}
	return c.JSON(http.StatusOK, response.SuccessfulOK())
}

func (h *Handler) DeleteBooks(c echo.Context) error {
	request := &googleCloud.DeleteUsersForm{}
	cc := c.(*middleware.CustomContext)
	if err := cc.BindAndValidate(request); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	err := h.service.DeleteBooks(c, request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	response := &mongodb.Response{}
	return c.JSON(http.StatusOK, response.SuccessfulOK())
}

func (h *Handler) UploadImage(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	request := &googleCloud.UploadForm{
		File: file,
	}

	url, err := h.service.UploadImage(c, request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"link": url,
	})
	
}
