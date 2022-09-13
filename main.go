package main

import (
	"context"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	//"github.com/khotchapan/KonLakRod-api/routes"
	//"github.com/khotchapan/KonLakRod-api/configs"
	"github.com/khotchapan/KonLakRod-api/connection"
	coreContext "github.com/khotchapan/KonLakRod-api/internal/core/context"
	coreValidator "github.com/khotchapan/KonLakRod-api/internal/core/validator"
	guestEndpoint "github.com/khotchapan/KonLakRod-api/internal/handlers/guest"
	googleCloud "github.com/khotchapan/KonLakRod-api/lagacy/google/google_cloud"
	users "github.com/khotchapan/KonLakRod-api/mongodb/user"
	"github.com/khotchapan/KonLakRod-api/services/test"
	"github.com/khotchapan/KonLakRod-api/services/user"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

func main() {
	var (
		e           = initEcho()
		dbMonggo, _ = newMongoDB()
		gcs         = googleCloud.NewGoogleCloudStorage(dbMonggo)
	)
	app := context.WithValue(context.Background(), connection.ConnectionInit,
		connection.Connection{
			Monggo: dbMonggo,
			GCS:    gcs,
		})
	collection := context.WithValue(context.Background(), connection.CollectionInit,
		connection.Collection{
			Users: users.NewRepo(dbMonggo),
		})
	//===============================================================================
	// Configure middleware with the custom claims type
	config := middleware.JWTConfig{
		//Claims:     &jwtCustomClaims{},
		Claims:     &coreContext.Claims{},
		SigningKey: []byte("secret"),
		//ContextKey: "user",
	}
	checkSessionMiddleware := middleware.JWTWithConfig(config)
	//===============================================================================
	//home
	e.GET(path.Join("/"), Version)
	api := e.Group("/v1/api")
	api.GET("", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// Login route
	//api.POST("/login", login)

	// Unauthenticated route
	api.GET("/", accessible)
	// Restricted group
	r := api.Group("/restricted", checkSessionMiddleware)
	{
		//api.Use(checkSessionMiddleware)
		r.GET("", restricted)
	}
	guest := guestEndpoint.NewHandler(guestEndpoint.NewService(app, collection))
	guestGroup := api.Group("/guest")
	{
		guestGroup.POST("/login", guest.PostLoginUsers)
		guestGroup.POST("/login-test", guest.Login)
	}

	//user
	users := user.NewHandler(user.NewService(app, collection))
	usersGroup := api.Group("/users")
	usersGroup.GET("", users.GetAllUsers)
	usersGroup.GET("/:id", users.GetOneUsers)
	usersGroup.POST("", users.PostUsers)
	usersGroup.PUT("", users.PutUsers)
	usersGroup.DELETE("", users.DeleteUsers)
	usersGroup.POST("/upload", users.UploadFile)
	usersGroup.POST("image/upload", users.UploadFileUsers)

	// test zone
	testService := test.NewHandler(test.NewService(app, collection))
	testGroup := api.Group("/tests")
	testGroup.GET("/google-cloud/books", testService.GetFile)
	testGroup.GET("/google-cloud/books/:id", testService.GetOneGoogleCloudBooks)
	testGroup.POST("/google-cloud/books", testService.PostGoogleCloudBooks)
	testGroup.PUT("/google-cloud/books", testService.PutBooks)
	testGroup.DELETE("/google-cloud/books", testService.DeleteBooks)
	testGroup.POST("/google-cloud/image/upload", testService.UploadImage)
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

	e.Logger.Fatal(e.Start(address))

}
func Version(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{"version": 1.9})
}

func initEcho() *echo.Echo {
	e := echo.New()
	e.Validator = coreValidator.NewValidator(validator.New())
	e.Use(coreContext.SetCustomContext)
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

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

// jwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}



func accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*coreContext.Claims)
	name := claims.Subject
	return c.String(http.StatusOK, "Welcome "+name+"!")
}
