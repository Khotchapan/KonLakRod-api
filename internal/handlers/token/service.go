package token

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/khotchapan/KonLakRod-api/internal/core/connection"
	"github.com/khotchapan/KonLakRod-api/internal/entities"
	coreMiddleware "github.com/khotchapan/KonLakRod-api/internal/middleware"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type ServiceInterface interface {
	Create(c echo.Context, u *entities.Users) (*entities.TokenResponse, error)
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

func (s *Service) Create(c echo.Context, u *entities.Users) (*entities.TokenResponse, error) {
	token, err := s.createJWTToken(c, u)
	if err != nil {
		return nil, err
	}

	return token, nil
}
func (s *Service) createJWTTokenTest(c echo.Context, u *entities.Users) (*entities.TokenResponse, error) {
	rto, err := s.createRefreshToken(u)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	claims := &coreMiddleware.Claims{}
	claims.Subject = "access_token"
	claims.Issuer = "KonLakRod"
	claims.IssuedAt = now.Unix()
	claims.ExpiresAt = now.Add(time.Hour * 24).Unix()
	claims.Id = uuid.New().String()
	claims.UserID = &u.ID
	claims.Roles = u.Roles
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return nil, err
	}
	td := &entities.TokenDetails{}
	td.AtExpires = time.Now().Add(time.Hour * 2).Unix()
	td.AccessUuid = uuid.New().String()
	td.RtExpires = time.Now().Add(time.Minute * 24 * 7).Unix()
	td.RefreshUuid = uuid.New().String()
	td.AccessToken, err = token.SignedString([]byte("secret"))
	if err != nil {
		return nil, err
	}
	//====================================================================
	rtToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := rtToken.Claims.(jwt.MapClaims)
	rtClaims["id"] = &u.ID
	rtClaims["sub"] = "refresh_token"
	rtClaims["exp"] = td.RtExpires
	rtClaims["jti"] = td.RefreshUuid
	td.RefreshToken, err = rtToken.SignedString([]byte("secret"))
	if err != nil {
		return nil, err
	}

	//====================================================================
	tk := &entities.Tokens{RefreshToken: td.RefreshToken,
		UserRefId: &u.ID}

	err = s.collection.Tokens.Create(tk)
	if err != nil {
		return nil, err
	}

	tkr := &entities.TokenResponse{
		AccessToken:      &t,
		RefreshToken:     &rto,
		// AccessTokenTest:  &td.AccessToken,
		// RefreshTokenTest: &td.RefreshToken,
	}

	return tkr, nil
}
func (s *Service) createJWTToken(c echo.Context, u *entities.Users) (*entities.TokenResponse, error) {
	now := time.Now()
	tokenDetailsTest := &entities.TokenDetailsTest{}
	tokenDetailsTest.IssuedAt = now.Unix()
	tokenDetailsTest.AccessTokenExpiresAt = now.Add(time.Hour * 1).Unix()
	tokenDetailsTest.RefreshTokenExpiresAt = time.Now().Add(time.Hour * 24 * 14).Unix()
	tokenDetailsTest.AccessTokenId = uuid.New().String()
	tokenDetailsTest.RefreshTokenId = uuid.New().String()

	claims := &coreMiddleware.Claims{}
	claims.Subject = "access_token"
	claims.Issuer = "KonLakRod"
	claims.IssuedAt = tokenDetailsTest.IssuedAt
	claims.ExpiresAt = tokenDetailsTest.AccessTokenExpiresAt
	claims.Id = tokenDetailsTest.AccessTokenId
	claims.UserID = &u.ID
	claims.Roles = u.Roles
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	accessToken, err := token.SignedString([]byte("secret"))
	if err != nil {
		return nil, err
	}
	tokenDetailsTest.AccessToken = accessToken

	//====================================================================
	rtToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := rtToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = "refresh_token"
	rtClaims["iss"] = "KonLakRod"
	rtClaims["iat"] = tokenDetailsTest.IssuedAt
	rtClaims["exp"] = tokenDetailsTest.RefreshTokenExpiresAt
	rtClaims["jti"] = tokenDetailsTest.RefreshTokenId
	//rtClaims["user_id"] = &u.ID
	//rtClaims["roles"] = u.Roles
	refreshToken, err := rtToken.SignedString([]byte("secret"))
	if err != nil {
		return nil, err
	}
	tokenDetailsTest.RefreshToken = refreshToken

	//====================================================================
	tk := &entities.Tokens{RefreshToken: tokenDetailsTest.RefreshToken,
		UserRefId: &u.ID}

	err = s.collection.Tokens.Create(tk)
	if err != nil {
		return nil, err
	}

	tkr := &entities.TokenResponse{
		AccessToken:      &tokenDetailsTest.AccessToken,
		RefreshToken:     &tokenDetailsTest.RefreshToken,
	}

	return tkr, nil
}

func (s *Service) RefreshToken(c echo.Context, request *RefreshTokenForm) (*entities.TokenResponse, error) {
	token, err := s.verifyToken(*request.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid token or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token or expired token")
		//return nil, echo.ErrUnauthorized
	}

	id, ok := claims["user_id"].(string)
	log.Println("user_id:", id)
	if !ok {
		return nil, errors.New("invalid JWT Payload")
	}
	//return nil, nil
	//============================================================================
	tk := &entities.Tokens{}
	log.Println("request.RefreshToken:", *request.RefreshToken)
	err = s.collection.Tokens.FindOneByRefreshToken(request.RefreshToken, tk)
	log.Println("tk:", tk)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil,errors.New("error no documents")
		}
		return nil,err
	}
	err = s.collection.Tokens.Delete(tk)
	if err != nil {
		return nil, err
	}
	us := &entities.Users{}
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

func (s *Service) createRefreshToken(u *entities.Users) (string, error) {
	rts := fmt.Sprintf("%d%s", u.ID, time.Now().String())
	h := sha1.New()
	_, err := h.Write([]byte(rts))
	if err != nil {
		return "", err
	}
	res := hex.EncodeToString(h.Sum(nil))
	//log.Println("EncodeToString:", res)
	return res, nil
}

func (s *Service) verifyToken(tokenStr string) (*jwt.Token, error) {
	// Parse takes the token string and a function for looking up the key.
	// The latter is especially useful if you use multiple keys for your application.
	// The standard is to use 'kid' in the head of the token to identify
	// which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("secret"), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
