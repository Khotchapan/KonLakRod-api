package token

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/khotchapan/KonLakRod-api/internal/core/connection"
	coreContext "github.com/khotchapan/KonLakRod-api/internal/core/context"
	"github.com/khotchapan/KonLakRod-api/internal/core/mongodb/user"
	"github.com/khotchapan/KonLakRod-api/internal/entities"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type ServiceInterface interface {
	Create(c echo.Context, u *user.Users) (*entities.TokenResponse, error)
	RefreshToken(c echo.Context, request *RefreshTokenForm) (*entities.TokenResponse, error)
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

func (s *Service) Create(c echo.Context, u *user.Users) (*entities.TokenResponse, error) {
	log.Println("========STEP4========")
	token, err := s.createJWTToken(c, u)
	if err != nil {
		return nil, err
	}

	return token, nil
}
func (s *Service) createJWTToken(c echo.Context, u *user.Users) (*entities.TokenResponse, error) {
	rto, err := s.createRefreshToken(u)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	claims := &coreContext.Claims{}
	claims.Subject = "access_token"
	claims.Issuer = "KonLakRod"
	claims.IssuedAt = now.Unix()
	claims.ExpiresAt = now.Add(time.Hour * 24).Unix()
	claims.Id = uuid.New().String()
	claims.UserID = &u.ID
	claims.Roles = u.Roles
	claims.User = u
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return nil, err
	}
	tk := &entities.Tokens{RefreshToken: rto,
		UserRefId: &u.ID}

	err = s.collection.Tokens.Create(tk)
	if err != nil {
		return nil, err
	}

	tkr := &entities.TokenResponse{
		AccessToken:  &t,
		RefreshToken: &rto,
	}
	return tkr, nil
}
func (s *Service) RefreshToken(c echo.Context, request *RefreshTokenForm) (*entities.TokenResponse, error) {
	// err := s.verifyToken(request.RefreshToken)
	// if err != nil {
	// 	return nil, errors.New("invalid token or expired token")
	// }
	tk := &entities.Tokens{}
	log.Println("request.RefreshToken:", *request.RefreshToken)
	err := s.collection.Tokens.FindOneByRefreshToken(request.RefreshToken, tk)
	log.Println("tk:", tk)
	if err != nil && err != mongo.ErrNoDocuments {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		return nil, errors.New("error no documents")
	}
	err = s.collection.Tokens.Delete(tk)
	if err != nil {
		return nil, err
	}
	us := &user.Users{}
	err = s.collection.Users.FindOneByObjectID(tk.UserRefId, us)
	if err != nil {
		return nil, err
	}
	log.Println("us:", us)
	t, err := s.Create(c, us)
	if err != nil {

		return nil, err
	}
	return t, nil

}

func (s *Service) createRefreshToken(u *user.Users) (string, error) {
	rts := fmt.Sprintf("%d%s", u.ID, time.Now().String())
	h := sha1.New()
	_, err := h.Write([]byte(rts))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func (s *Service) verifyToken(tokenStr string) error {
	// Parse takes the token string and a function for looking up the key.
	// The latter is especially useful if you use multiple keys for your application.
	// The standard is to use 'kid' in the head of the token to identify
	// which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	_, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("secret"), nil
	})

	if err != nil {
		return err
	}

	return nil
}
