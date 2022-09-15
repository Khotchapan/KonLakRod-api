package token

import (
	"context"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/khotchapan/KonLakRod-api/internal/core/connection"
	coreContext "github.com/khotchapan/KonLakRod-api/internal/core/context"
	"github.com/khotchapan/KonLakRod-api/internal/core/mongodb/user"
	"github.com/labstack/echo/v4"
)

type ServiceInterface interface {
	Create2(c echo.Context, u *user.Users) (*string, error)
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

func (s *Service) Create2(c echo.Context, u *user.Users) (*string, error) {
	log.Println("========STEP4========")
	token, err := s.createJWTToken2(c, u)
	if err != nil {
		return nil, err
	}

	return token, nil
}
func (s *Service) createJWTToken2(c echo.Context, u *user.Users) (*string, error) {
	now := time.Now()
	claims := &coreContext.Claims{}
	claims.Subject = "access_token"
	claims.Issuer = "KonLakRod"
	claims.IssuedAt = now.Unix()
	claims.ExpiresAt = now.Add(time.Hour * 24).Unix()
	claims.Id = uuid.New().String()
	claims.UserID = u.ID.Hex()
	claims.Roles = u.Roles
	claims.User = u
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return nil, err
	}
	return &t, nil
}
