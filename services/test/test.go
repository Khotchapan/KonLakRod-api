package test

import (
	"context"
	"github.com/khotchapan/KonLakRod-api/connection"
	googleCloud "github.com/khotchapan/KonLakRod-api/lagacy/google/google_cloud"
	"github.com/labstack/echo/v4"
)

type TestInterface interface {
	FindAllBooks(c echo.Context) ([]*googleCloud.Books, error)
	FindOneBooks(c echo.Context, request *GetOneGoogleCloudBooksForm) ([]*googleCloud.Books, error)
	CreateBooks(c echo.Context, request *googleCloud.CreateBooksForm) error
	UpdateBooks(c echo.Context, request *googleCloud.UpdateBooksForm) error
	DeleteBooks(c echo.Context, request *googleCloud.DeleteUsersForm) error
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

func (s *Service) FindAllBooks(c echo.Context) ([]*googleCloud.Books, error) {
	response, err := s.con.GCS.FindAllBooks()
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *Service) FindOneBooks(c echo.Context, request *GetOneGoogleCloudBooksForm) ([]*googleCloud.Books, error) {
	response, err := s.con.GCS.FindOneBooks(request.ID)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *Service) CreateBooks(c echo.Context, request *googleCloud.CreateBooksForm) error {
	books := &googleCloud.Books{}
	b := request.Fill(books)
	err := s.con.GCS.CreateBooks(b)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateBooks(c echo.Context, request *googleCloud.UpdateBooksForm) error {
	books := &googleCloud.Books{}
	b := request.Fill(books)
	err := s.con.GCS.UpdateBooks(b)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) DeleteBooks(c echo.Context, request *googleCloud.DeleteUsersForm) error {
	err := s.con.GCS.DeleteBooks(request)
	if err != nil {
		return err
	}
	return nil
}