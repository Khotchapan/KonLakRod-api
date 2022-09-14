package connection

import (
	"context"
	"log"
	"github.com/khotchapan/KonLakRod-api/internal/core/mongodb/user"
	googleCloud "github.com/khotchapan/KonLakRod-api/internal/lagacy/google/google_cloud"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	ConnectionInit = "S29uTGFrUm9kLWFwaS1jb25uZWN0"
	CollectionInit = "S29uTGFrUm9kLWNvbm5lY3QtY29sbGVjdGlvbg=="
)

type Connection struct {
	Monggo *mongo.Database
	// AzureRoom       azure.AzureInterface
	// PhramcareClient phramcareClient.PhramcareSendInterface
	GCS googleCloud.IGCS
}

type Collection struct {
	Users user.UsersInterface
	//TokenService token.ServiceInterface
}

func GetConnect(ctx context.Context, k string) *Connection {
	if v, ok := ctx.Value(k).(Connection); ok {
		return &v
	}
	log.Panicln("Service can't create Connection or ctx not match")
	return nil
}

func GetCollection(ctx context.Context, k string) *Collection {
	if v, ok := ctx.Value(k).(Collection); ok {
		return &v
	}
	log.Panicln("Seveice can't create Collection or ctx not math")
	return nil
}
