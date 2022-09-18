package guest

import (
	"context"
	"errors"
	"log"
	"github.com/khotchapan/KonLakRod-api/internal/core/bcrypt"
	"github.com/khotchapan/KonLakRod-api/internal/core/connection"
	"github.com/khotchapan/KonLakRod-api/internal/entities"
	"github.com/khotchapan/KonLakRod-api/internal/handlers/token"
	"github.com/labstack/echo/v4"
)

type ServiceInterface interface {
	LoginUsers(c echo.Context, request *LoginUsersForm) (*entities.TokenResponse, error)
}

type Service struct {
	con        *connection.Connection
	collection *connection.Collection
	//tokenService *token.Service
	tokenService token.ServiceInterface
}

func NewService(app, collection context.Context) *Service {
	return &Service{
		con:          connection.GetConnect(app, connection.ConnectionInit),
		collection:   connection.GetCollection(collection, connection.CollectionInit),
		tokenService: token.NewService(app, collection),
	}
}

func (s *Service) LoginUsers(c echo.Context, request *LoginUsersForm) (*entities.TokenResponse, error) {
	us := &entities.Users{}
	err := s.collection.Users.FindOneByName(request.Username, us)
	if err != nil {
		return nil, err
	}
	if !bcrypt.ComparePassword(*request.Password, us.PasswordHash) {
		log.Println("check")
		return nil, errors.New("password is incorrect")
	}
	token, err := s.tokenService.Create(c, us)
	if err != nil {
		return nil, err
	}

	return token, nil
}
