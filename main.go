package main

import (
	"context"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"github.com/khotchapan/KonLakRod-api/internal/core/connection"
	coreMiddleware "github.com/khotchapan/KonLakRod-api/internal/middleware"
	users "github.com/khotchapan/KonLakRod-api/internal/core/mongodb/user"
	tokens "github.com/khotchapan/KonLakRod-api/internal/core/mongodb/token"
	coreValidator "github.com/khotchapan/KonLakRod-api/internal/core/validator"
	googleCloud "github.com/khotchapan/KonLakRod-api/internal/lagacy/google/google_cloud"
	"github.com/khotchapan/KonLakRod-api/internal/router"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
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
			Tokens: tokens.NewRepo(dbMonggo),
		})
	options := &router.Options{
		App:        app,
		Collection: collection,
		Echo:       e,
	}
	router.Router(options)

	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = "80" // Default port if not specified
	}
	//========================================================

	//========================================================
	address := fmt.Sprintf("%s:%s", "0.0.0.0", port)
	fmt.Println("address:", address)
	e.Logger.Fatal(e.Start(address))

}

func initEcho() *echo.Echo {
	e := echo.New()
	// e.HideBanner = false
	// e.HidePort = false
	// e.Debug = false
	// e.HideBanner = true
	e.Validator = coreValidator.NewValidator(validator.New())
	// Middleware
	e.Use(coreMiddleware.SetCustomContext)
	e.Use(middleware.Logger())    // Log everything to stdout
	e.Use(middleware.Recover())   // Recover from all panics to always have your server up
	e.Use(middleware.RequestID()) // Generate a request id on the HTTP response headers for identification
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
