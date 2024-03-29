package connection

import (
	"context"
	"log"
	"github.com/khotchapan/KonLakRod-api/internal/core/memory"
	postReply "github.com/khotchapan/KonLakRod-api/internal/core/mongodb/post_reply"
	postTopic "github.com/khotchapan/KonLakRod-api/internal/core/mongodb/post_topic"
	"github.com/khotchapan/KonLakRod-api/internal/core/mongodb/token"
	"github.com/khotchapan/KonLakRod-api/internal/core/mongodb/user"
	googleCloud "github.com/khotchapan/KonLakRod-api/internal/lagacy/google/google_cloud"
	"go.mongodb.org/mongo-driver/mongo"
)
type contextKey string
const (
	ConnectionInit contextKey = contextKey("S29uTGFrUm9kLWFwaS1jb25uZWN0")
	CollectionInit contextKey = contextKey("S29uTGFrUm9kLWNvbm5lY3QtY29sbGVjdGlvbg==")
)

type Connection struct {
	Mongo *mongo.Database
	GCS   googleCloud.GoogleCloudInterface
	Redis *memory.Redis
}
type Collection struct {
	Users     user.RepoInterface
	Tokens    token.RepoInterface
	PostTopic postTopic.RepoInterface
	PostReply postReply.RepoInterface
}

func GetConnect(ctx context.Context, k contextKey) *Connection {

	if v, ok := ctx.Value(k).(Connection); ok {
		//log.Println("GetConnect:",v)
		return &v
	}
	log.Panicln("Service can't create Connection or ctx not match")
	return nil
}

func GetCollection(ctx context.Context, k contextKey) *Collection {
	if v, ok := ctx.Value(k).(Collection); ok {
		//log.Println("GetCollection:",v)
		return &v
	}
	log.Panicln("Seveice can't create Collection or ctx not math")
	return nil
}
