package test

import (
	"context"

	"github.com/khotchapan/KonLakRod-api/connection"
	googleCloud "github.com/khotchapan/KonLakRod-api/lagacy/google/google_cloud"
	"github.com/labstack/echo/v4"
)

type TestInterface interface {
	FindAllBooks(c echo.Context) ([]*googleCloud.Book, error)
	FindOneBooks(c echo.Context,request *GetOneGoogleCloudForm) ([]*googleCloud.Book, error)
}

type Service struct {
	con        *connection.Connection
	collection *connection.Collection
}

func NewService(app, collection context.Context) *Service {
	return &Service{
		con:        connection.GetConnect(app, connection.ConnectionInit),
		collection: connection.GetCollection(collection, connection.CollectionInit),
	}
}
func (s *Service) FindAllBooks(c echo.Context) ([]*googleCloud.Book, error) {
	response, err := s.con.GCS.FindAllBooks()
	if err != nil {
		return nil, err
	}

	return response, nil
}
func (s *Service) FindOneBooks(c echo.Context,request *GetOneGoogleCloudForm) ([]*googleCloud.Book, error) {
	response, err := s.con.GCS.FindOneBooks(request.ID)
	if err != nil {
		return nil, err
	}

	return response, nil
}

