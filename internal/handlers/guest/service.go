package guest

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

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
	con          *connection.Connection
	collection   *connection.Collection
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
	err := s.collection.Users.FindOneByUserName(request.Username, us)
	if err != nil {
		return nil, err
	}
	if !bcrypt.ComparePassword(*request.Password, us.PasswordHash) {
		return nil, errors.New("password is incorrect")
	}
	tokenDetails, err := s.tokenService.GenerateTokensAndSetDatabase(c, us)
	if err != nil {
		return nil, err
	}
	//s.con.GCS.CreateBooks()
	token := &entities.TokenResponse{
		AccessToken:      &tokenDetails.AccessToken,
		RefreshToken:     &tokenDetails.RefreshToken,
	}
	json, err := json.Marshal(map[string]string{"some": "value"})
	if err != nil {
		return nil, err
	}
	s.con.Redis.Test("name", json, (20)*time.Second)
	//s.con.Redis.SetKey("name", json, (10)*time.Second)
	log.Println("=========================================")
	return token, nil
}
