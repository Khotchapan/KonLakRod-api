package main

import (
	"context"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"github.com/khotchapan/KonLakRod-api/connection"
	coreContext "github.com/khotchapan/KonLakRod-api/internal/core/context"
	users "github.com/khotchapan/KonLakRod-api/internal/core/mongodb/user"
	coreValidator "github.com/khotchapan/KonLakRod-api/internal/core/validator"
	"github.com/khotchapan/KonLakRod-api/internal/router"
	googleCloud "github.com/khotchapan/KonLakRod-api/lagacy/google/google_cloud"
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
		})
	options := &router.Options{
		App:        app,
		Collection: collection,
		Echo:       e,
	}
	router.Router(options)

	godotenv.Load()
	port := os.Getenv("PORT")
	port = "1323"
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
	e.Use(coreContext.SetCustomContext)
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
