package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/khotchapan/KonLakRod-api/internal/core/bcrypt"
	"github.com/khotchapan/KonLakRod-api/internal/core/connection"
	coreContext "github.com/khotchapan/KonLakRod-api/internal/core/context"
	"github.com/khotchapan/KonLakRod-api/internal/core/mongodb"
	"github.com/khotchapan/KonLakRod-api/internal/core/mongodb/user"
	"github.com/khotchapan/KonLakRod-api/internal/entities"
	googleCloud "github.com/khotchapan/KonLakRod-api/internal/lagacy/google/google_cloud"
	"github.com/labstack/echo/v4"
)

type UserInterface interface {
	CallGetMe(c echo.Context) (*user.Users, error)
	FindAllUsers(c echo.Context, request *user.GetAllUsersForm) (*mongodb.Page, error)
	FindOneUsers(c echo.Context, request *GetOneUsersForm) (*user.Users, error)
	CreateUsers(c echo.Context, request *CreateUsersForm) error
	UpdateUsers(c echo.Context, request *UpdateUsersForm) error
	DeleteUsers(c echo.Context, request *DeleteUsersForm) error
	UploadFile(c echo.Context, req UploadForm) (string, error)
	UploadFileUsers(c echo.Context, req *googleCloud.UploadForm) (*googleCloud.ImageStructure, error)
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
func (s *Service) CallGetMe(c echo.Context) (*user.Users, error) {
	response := &user.Users{}
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*coreContext.Claims)
	log.Println("UserID:", claims.UserID)
	userID := claims.UserID
	err := s.collection.Users.FindOneByObjectID(userID, response)
	if err != nil {
		return nil, err
	}
	//err := s.collection.Users.FindOneByID(c.Get("user").(*jwt.User).ID.Hex(), response)
	// if err != nil {
	// 	return nil, errs.New(http.StatusConflict, "10003", "user not found")
	// }

	if response.Image != "" && !strings.Contains(response.Image, "http") {
		url, err := s.con.GCS.SignedURL(response.Image)
		if err != nil {
			return nil, errors.New("can not singed url")
		}
		response.Image = url
	}

	// if response.HealthInfo.Birthday != "" {
	// 	birthday, err := time.Parse("2006-01-02", response.HealthInfo.Birthday)
	// 	if err != nil {
	// 		return nil, errs.NewBadRequest("Can not convert to time", err.Error())
	// 	}

	// 	y, m, d := birthday.Date()
	// 	dob := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	// 	response.HealthInfo.Age = age.Age(dob)
	// }

	return response, nil
}
func (s *Service) FindAllUsers(c echo.Context, request *user.GetAllUsersForm) (*mongodb.Page, error) {
	//objectUserID := &c.Get("user").(*jwt.User).ID
	response, err := s.collection.Users.FindAllUsers(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}
func (s *Service) FindOneUsers(c echo.Context, request *GetOneUsersForm) (*user.Users, error) {
	response := &user.Users{}
	err := s.collection.Users.FindOneByObjectID(request.ID, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *Service) CreateUsers(c echo.Context, request *CreateUsersForm) error {
	us := &user.Users{}
	password, err := bcrypt.GeneratePassword(*request.Password)
	if err != nil {
		//c.Error(err)
		return err
	}
	us.PasswordHash = password
	us.Roles = []string{entities.UserRole}
	u := request.fill(us)
	err = s.collection.Users.Create(u)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateUsers(c echo.Context, request *UpdateUsersForm) error {
	us := &user.Users{}
	err := s.collection.Users.FindOneByObjectID(request.ID, us)
	if err != nil {
		return err
	}
	u := request.fill(us)
	err = s.collection.Users.Update(u)
	if err != nil {
		return err
	}
	return nil
}
func (s *Service) DeleteUsers(c echo.Context, request *DeleteUsersForm) error {
	u := &user.Users{
		Model: mongodb.Model{ID: *request.ID},
	}
	err := s.collection.Users.Delete(u)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UploadFile(c echo.Context, req UploadForm) (string, error) {
	src, err := req.File.Open()
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("test/%s.png", uuid.New().String())

	obj, _ := s.con.GCS.UploadFilePrivate(src, path)
	return s.con.GCS.SignedURL(obj)
}

func (s *Service) UploadFileUsers(c echo.Context, request *googleCloud.UploadForm) (*googleCloud.ImageStructure, error) {
	imageStructure, err := s.con.GCS.UploadFileUsers(request)
	if err != nil {
		return nil, err
	}
	//log.Println("objectName:", *objectName)
	signedUrl, err := s.con.GCS.SignedURL(imageStructure.ImageName)
	if err != nil {
		return nil, err
	}
	imageStructure.URL = signedUrl
	return imageStructure, nil
}
