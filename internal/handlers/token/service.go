package token

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/khotchapan/KonLakRod-api/internal/core/connection"
	coreContext "github.com/khotchapan/KonLakRod-api/internal/core/context"
	"github.com/khotchapan/KonLakRod-api/internal/core/mongodb/user"
	"github.com/khotchapan/KonLakRod-api/internal/entities"
	"github.com/labstack/echo/v4"
)

type ServiceInterface interface {
	Create(c echo.Context, u *user.Users) (*entities.Token, error)
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
	claims.Issuer = "test.plus"
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
	privateKey, _ := hex.DecodeString(RsaPrivateKey)
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		log.Println(err)
	}
	//========================================
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	tokenString, err := token.SignedString(signKey)
	if err != nil {
		log.Println("tokenString")
		return nil, err
	}
	t.Token = tokenString
	return t, nil
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
	claims.Subject = fmt.Sprint("Hello World1234")
	claims.Issuer = "kidscare.plus"
	claims.IssuedAt = now.Unix()
	claims.ExpiresAt = now.Add(time.Hour * 24).Unix()
	claims.Roles = u.Roles
	claims.UserID = u.ID.Hex()
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
