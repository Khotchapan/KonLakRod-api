package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
	"github.com/go-playground/validator"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/khotchapan/KonLakRod-api/internal/core/connection"
	"github.com/khotchapan/KonLakRod-api/internal/core/memory"
	postReply "github.com/khotchapan/KonLakRod-api/internal/core/mongodb/post_reply"
	postTopic "github.com/khotchapan/KonLakRod-api/internal/core/mongodb/post_topic"
	tokens "github.com/khotchapan/KonLakRod-api/internal/core/mongodb/token"
	users "github.com/khotchapan/KonLakRod-api/internal/core/mongodb/user"
	"github.com/khotchapan/KonLakRod-api/internal/core/utils"
	coreValidator "github.com/khotchapan/KonLakRod-api/internal/core/validator"
	googleCloud "github.com/khotchapan/KonLakRod-api/internal/lagacy/google/google_cloud"
	coreMiddleware "github.com/khotchapan/KonLakRod-api/internal/middleware"
	"github.com/khotchapan/KonLakRod-api/internal/router"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//var ctx context.Context = context.Background()

func initViper() {

	viper.AddConfigPath("configs")                         // ระบุ path ของ config file
	viper.SetConfigName("config")                          // ชื่อ config file
	viper.AutomaticEnv()                                   // อ่าน value จาก ENV variable
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // แปลง _ underscore ใน env เป็น . dot notation ใน viper
	// read config
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("cannot read in viper config:%s", err)
	}
	log.Println(viper.Get("app.env"))
}
func initGoDotEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
func init() {
	log.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
	initGoDotEnv()
	initViper()
	log.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
}
func main() {
	var (
		e             = initEcho()
		mongodb, _    = newMongoDB()
		redisDatabase = newRedis()
		gcs           = googleCloud.NewGoogleCloudStorage(mongodb)
	)
	app := context.WithValue(context.Background(), connection.ConnectionInit,
		connection.Connection{
			Mongo: mongodb,
			GCS:   gcs,
			Redis: memory.New(redisDatabase),
		})
	collection := context.WithValue(context.Background(), connection.CollectionInit,
		connection.Collection{
			Users:     users.NewRepo(mongodb),
			Tokens:    tokens.NewRepo(mongodb),
			PostTopic: postTopic.NewRepo(mongodb),
			PostReply: postReply.NewRepo(mongodb),
		})
	options := &router.Options{
		App:        app,
		Collection: collection,
		Echo:       e,
	}
	router.Router(options)
	port := utils.Getenv("PORT", viper.GetString("app.port"))
	address := fmt.Sprintf("%s:%s", "0.0.0.0", port)
	log.Println("address:", address)
	e.Logger.Fatal(e.Start(address))
}

func initEcho() *echo.Echo {
	e := echo.New()
	// e.HideBanner = false
	// e.HidePort = false
	// e.Debug = false
	// e.HideBanner = true
	//Validator
	e.Validator = coreValidator.NewValidator(validator.New())
	// Middleware
	e.Use(coreMiddleware.SetCustomContext)
	e.Use(middleware.Logger())    // Log everything to stdout
	e.Use(middleware.Recover())   // Recover from all panics to always have your server up
	e.Use(middleware.RequestID()) // Generate a request id on the HTTP response headers for identification
	return e
}

func newMongoDB() (*mongo.Database, *context.Context) {
	//EnvMongoURI := os.Getenv("MONGOURI")
	//viper.SetDefault("MONGO.HOST", "mongodb+srv://admin:1234@cluster0.6phd9zm.mongodb.net/?retryWrites=true&w=majority")
	EnvMongoURI := viper.GetString("MONGO.HOST")
	log.Println("EnvMongoURI", EnvMongoURI)
	client, err := mongo.NewClient(options.Client().ApplyURI(EnvMongoURI))
	if err != nil {
		log.Fatal(err)
	}

	contextDatabase, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(contextDatabase)
	if err != nil {
		log.Fatal(err)
	}

	//ping the database
	err = client.Ping(contextDatabase, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	return client.Database("konlakrod"), &contextDatabase
}

func newRedis() *redis.Client {
	host := utils.Getenv("REDIS_URI", "localhost")
	log.Println("HOST::::::::::", host)
	val := fmt.Sprintf("%s:%s", host, "6379")
	rdb := redis.NewClient(&redis.Options{
		//Addr:     "localhost:6379",
		Addr:     val,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("redis error:", err)
	}
	log.Println("redis:", pong)
	return rdb
}
