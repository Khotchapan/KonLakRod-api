package user

import (
	"context"
	"github.com/khotchapan/KonLakRod-api/connection"
	"github.com/khotchapan/KonLakRod-api/mongodb"
	"github.com/khotchapan/KonLakRod-api/mongodb/user"
	"github.com/labstack/echo/v4"
)

type UserInterface interface {
	FindAllUsers(c echo.Context, request *user.GetAllUsersForm) (*mongodb.Page, error)
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

func (s *Service) FindAllUsers(c echo.Context, request *user.GetAllUsersForm) (*mongodb.Page, error) {
	//objectUserID := &c.Get("user").(*jwt.User).ID
	response, err := s.collection.Users.FindAllUsers(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}
