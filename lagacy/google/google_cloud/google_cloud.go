package googleCloud

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type IGCS interface {
	UploadFile(file multipart.File, path string) (string, error)
	GetBucketName() string
	UploadFilePrivate(file multipart.File, path string) (string, error)
	SignedURL(object string) (string, error)
	FindAllBooks() ([]*Books, error)
	FindOneBooks(id *string) ([]*Books, error)
	CreateBooks(i *Books) error
}

type GoogleCloudStorage struct {
	cl         *storage.Client
	app        *firebase.App
	Client     *firestore.Client
	projectID  string
	bucketName string
	basePath   string
}

func NewGoogleCloudStorage(db *mongo.Database) IGCS {
	opt := option.WithCredentialsFile("internal/env/firebase_secret_key.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil
	}
	client, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	return &GoogleCloudStorage{
		Client: client,
	}
}

func (g *GoogleCloudStorage) UploadFile(file multipart.File, path string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	obj := g.basePath + path

	wc := g.cl.Bucket(g.bucketName).Object(obj).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("io.Copy: %v", err)
	}

	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %v", err)
	}

	err := g.cl.Bucket(g.bucketName).Object(obj).ACL().Set(ctx, storage.AllUsers, storage.RoleReader)
	if err != nil {
		return "", err
	}

	rObj := g.cl.Bucket(g.bucketName).Object(obj)
	return fmt.Sprintf("%s/%s/%s", viper.GetString("gcs.baseURL"), rObj.BucketName(), rObj.ObjectName()), nil
}

func (g *GoogleCloudStorage) UploadFilePrivate(file multipart.File, path string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	obj := g.basePath + path

	wc := g.cl.Bucket(g.bucketName).Object(obj).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("io.Copy: %v", err)
	}

	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %v", err)
	}

	rObj := g.cl.Bucket(g.bucketName).Object(obj)
	return rObj.ObjectName(), nil
}

func (g *GoogleCloudStorage) GetBucketName() string {
	return g.bucketName
}

func (g *GoogleCloudStorage) SignedURL(object string) (string, error) {
	// storage.SignedURL(g.bucketName, object, &storage.SignedURLOptions{
	// 	GoogleAccessID: ,
	// })
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	// defer cancel()
	// acc, _ := g.cl.ServiceAccount(ctx, g.projectID)
	return g.cl.Bucket(g.bucketName).SignedURL(object, &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(30 * time.Minute),
	})
}

func (g *GoogleCloudStorage) FindAllBooks() ([]*Books, error) {
	BooksData := []*Books{}
	iter := g.Client.Collection("books").Documents(context.Background())
	for {
		BookData := &Books{}
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, errors.New("Something wrong, please try again.")
		}
		mapstructure.Decode(doc.Data(), &BookData)
		BooksData = append(BooksData, BookData)
	}
	return BooksData, nil
}

func (g *GoogleCloudStorage) FindOneBooks(id *string) ([]*Books, error) {
	BooksData := []*Books{}
	iter := g.Client.Collection("books").Where("id", "==", id).Documents(context.Background())
	for {
		BookData := &Books{}
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		mapstructure.Decode(doc.Data(), &BookData)
		BooksData = append(BooksData, BookData)
	}
	return BooksData, nil
}
func (g *GoogleCloudStorage) CreateBooks(i *Books) error {
	// iter := g.Client.Collection("books").Where("id", "==", id).Documents(context.Background())
	// for {
	// 	BookData := &Books{}
	// 	doc, err := iter.Next()
	// 	if err == iterator.Done {
	// 		break
	// 	}
	// 	if err != nil {
	// 		log.Fatalf("Failed to iterate: %v", err)
	// 	}
	// 	mapstructure.Decode(doc.Data(), &BookData)
	// 	BooksData = append(BooksData, BookData)
	// }
	// return BooksData, nil
	uid := uuid.New()
	log.Println("uid:", uid)
	splitID := strings.Split(uid.String(), "-")
	log.Println("splitID:", splitID)
	id := splitID[0] + splitID[1] + splitID[2] + splitID[3] + splitID[4]
	log.Println("id:", id)
	i.ID = id
	_, _, err := g.Client.Collection("books").Add(context.Background(), i)
	if err != nil {
		log.Printf("An error has occurred: %s", err)
	}
	return nil
}
