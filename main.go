package main

import (
	"context"
	"fmt"

	"github.com/go-playground/validator"
	"github.com/joho/godotenv"

	//"github.com/khotchapan/KonLakRod-api/configs"
	"github.com/khotchapan/KonLakRod-api/connection"
	coreContext "github.com/khotchapan/KonLakRod-api/internal/core/context"
	coreValidator "github.com/khotchapan/KonLakRod-api/internal/core/validator"
	users "github.com/khotchapan/KonLakRod-api/mongodb/user"

	//"github.com/khotchapan/KonLakRod-api/routes"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/khotchapan/KonLakRod-api/services/user"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	var (
		e           = initEcho()
		dbMonggo, _ = newMongoDB()
	)
	app := context.WithValue(context.Background(), connection.ConnectionInit,
		connection.Connection{
			Monggo: dbMonggo,
		})
	collection := context.WithValue(context.Background(), connection.CollectionInit,
		connection.Collection{
			Users: users.NewRepo(dbMonggo),
		})
	e.GET(path.Join("/"), Version)
	api := e.Group("/v1/api")
	api.GET("", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	//user
	users := user.NewHandler(user.NewService(app, collection))
	usersGroup := api.Group("/users")
	usersGroup.GET("", users.GetAllUsers)
	usersGroup.GET("/:id", users.GetOneUsers)
	usersGroup.POST("", users.PostUsers)
	usersGroup.PUT("", users.PutUsers)
	usersGroup.DELETE("", users.DeleteUsers)
	godotenv.Load()
	port := os.Getenv("PORT")
	port = "1323"
	//========================================================

	//========================================================
	address := fmt.Sprintf("%s:%s", "0.0.0.0", port)
	fmt.Println("address:", address)
	//run database
	//configs.ConnectDB()
	//routes // test zone
	//routes.UserRoute(e) //add this
	e.Start(address)

}
func Version(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{"version": 1.2})
}

func initEcho() *echo.Echo {
	e := echo.New()
	e.Validator = coreValidator.NewValidator(validator.New())
	e.Use(coreContext.SetCustomContext)

	return e
}

func newMongoDB() (*mongo.Database, context.Context) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	EnvMongoURI := os.Getenv("MONGOURI")
	fmt.Println("EnvMongoURI", EnvMongoURI)
	client, err := mongo.NewClient(options.Client().ApplyURI(EnvMongoURI))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")
	//return client
	return client.Database("konlakrod"), ctx
}
