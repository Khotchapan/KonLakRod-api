package token

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/khotchapan/KonLakRod-api/connection"
	coreContext "github.com/khotchapan/KonLakRod-api/internal/core/context"
	"github.com/khotchapan/KonLakRod-api/internal/entities"
	"github.com/khotchapan/KonLakRod-api/mongodb/user"
	"github.com/labstack/echo/v4"
	"log"
	"time"
)

type ServiceInterface interface {
	Create(c echo.Context, u *user.Users) (*entities.Token, error)
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

func (s *Service) Create(c echo.Context, u *user.Users) (*entities.Token, error) {
	log.Println("========STEP4========")
	token, err := s.createJWTToken(c, u)
	if err != nil {
		return nil, err
	}

	return token, nil
}
func (s *Service) createJWTToken(c echo.Context, u *user.Users) (*entities.Token, error) {
	log.Println("========STEP5========")
	t := &entities.Token{
		UserID: u.ID.Hex(),
		User:   u,
	}
	// rto, err := s.createRefreshToken(u)
	// if err != nil {
	// 	c.Log.Error(err)
	// 	return nil, err
	// }
	// t.RefreshToken = rto

	claims := &coreContext.Claims{}
	now := time.Now()
	claims.Subject = fmt.Sprint(u.ID)
	claims.Issuer = "kidscare.plus"
	claims.IssuedAt = now.Unix()
	claims.ExpiresAt = now.Add(24 * time.Hour).Unix()
	claims.Roles = u.Roles
	// err = s.tokenRepo.Create(c.Db, t)
	// if err != nil {
	// 	c.Log.Error(err)
	// 	return nil, err
	// }
	//claims.RefreshTokenID = t.ID
	//========================================

	//========================================
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		log.Println("tokenString")
		return nil, err
	}
	t.Token = tokenString
	return t, nil
}
