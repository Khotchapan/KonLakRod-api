package user

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/khotchapan/KonLakRod-api/connection"
	googleCloud "github.com/khotchapan/KonLakRod-api/lagacy/google/google_cloud"
	"github.com/khotchapan/KonLakRod-api/mongodb"
	"github.com/khotchapan/KonLakRod-api/mongodb/user"
	"github.com/labstack/echo/v4"
)

type UserInterface interface {
	FindAllUsers(c echo.Context, request *user.GetAllUsersForm) (*mongodb.Page, error)
	FindOneUsers(c echo.Context, request *GetOneUsersForm) (*user.Users, error)
	CreateUsers(c echo.Context, request *CreateUsersForm) error
	UpdateUsers(c echo.Context, request *UpdateUsersForm) error
	DeleteUsers(c echo.Context, request *DeleteUsersForm) error
	UploadFile(c echo.Context, req UploadForm) (string, error)
	FindAllFile(c echo.Context) ([]*googleCloud.Book, error)
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
func (s *Service) FindOneUsers(c echo.Context, request *GetOneUsersForm) (*user.Users, error) {
	response := &user.Users{}
	err := s.collection.Users.FindOneUsers(request.ID, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
func (s *Service) CreateUsers(c echo.Context, request *CreateUsersForm) error {
	us := &user.Users{}
	//data := []*user.Users{}
	// err := s.collection.Users.Create(request, &data)
	// if err != nil {
	// 	return err
	// }
	u := request.fill(us)
	// if len(data) > 0 {
	// 	dm.PharmacyCode = data[0].PharmacyCode
	// }
	err := s.collection.Users.Create(u)
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
	dm := &user.Users{
		Model: mongodb.Model{ID: *request.ID},
	}
	err := s.collection.Users.Delete(dm)
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

func (s *Service) FindAllFile(c echo.Context) ([]*googleCloud.Book, error) {
	//response := []*googleCloud.Book{}

	data, err := s.con.GCS.FindAllFile()
	if err != nil {
		return nil, err
	}

	return data, nil
}
