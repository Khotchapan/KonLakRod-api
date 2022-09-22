package post_topic

import (
	"context"

	"github.com/khotchapan/KonLakRod-api/internal/core/connection"
	"github.com/khotchapan/KonLakRod-api/internal/core/mongodb"
	postTopic "github.com/khotchapan/KonLakRod-api/internal/core/mongodb/post_topic"
	"github.com/khotchapan/KonLakRod-api/internal/entities"
	"github.com/labstack/echo/v4"
)

type ServiceInterface interface {
	CreatePostTopic(c echo.Context, request *CreatePostTopicForm) error
	UpdatePostTopic(c echo.Context, request *UpdatePostTopicForm) error
	DeletePostTopic(c echo.Context, request *DeletePostTopicForm) error
	FindOnePostTopic(c echo.Context, request *GetOneTopicForm) (*entities.PostTopic, error)
	FindAllPostTopic(c echo.Context, request *postTopic.GetAllPostTopicForm) (*mongodb.Page, error)
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

func (s *Service) CreatePostTopic(c echo.Context, request *CreatePostTopicForm) error {
	us := &entities.PostTopic{}
	u := request.fill(us)
	err := s.collection.PostTopic.Create(u)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdatePostTopic(c echo.Context, request *UpdatePostTopicForm) error {
	books := &entities.PostTopic{}
	b := request.Fill(books)
	err := s.collection.PostTopic.Update(b)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) DeletePostTopic(c echo.Context, request *DeletePostTopicForm) error {
	err := s.collection.PostTopic.Delete(request)
	if err != nil {
		return err
	}
	return nil
}
func (s *Service) FindOnePostTopic(c echo.Context, request *GetOneTopicForm) (*entities.PostTopic, error) {
	response := &entities.PostTopic{}
	err := s.collection.PostTopic.FindOneByObjectID(request.ID, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *Service) FindAllPostTopic(c echo.Context, request *postTopic.GetAllPostTopicForm) (*mongodb.Page, error) {
	// cc := c.(*middleware.CustomContext)
	// userID := cc.GetClaims().UserID
	// log.Println("userID:", userID)
	// user := c.Get("user").(*jwt.Token)
	// claims := user.Claims.(*coreMiddleware.Claims)
	// log.Println("claims.UserID:", claims.UserID)

	response, err := s.collection.PostTopic.FindAllPostTopic(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}
