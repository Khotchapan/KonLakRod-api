package connection

import (
	"context"
	"encoding/base64"
	"github.com/khotchapan/KonLakRod-api/mongodb/user"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)
const (
	ConnectionInit = "S29uTGFrUm9kLWFwaS1jb25uZWN0"
	CollectionInit = "S29uTGFrUm9kLWNvbm5lY3QtY29sbGVjdGlvbg=="
)

type Connection struct {
	Monggo *mongo.Database
	// Redis           *memstore.Redis
	// Secret          secret.Secret
	// SMS             ant.SMSer
	// AzureRoom       azure.AzureInterface
	// PhramcareClient phramcareClient.PhramcareSendInterface
	// GG              google.SigninInterface
	// JWT             *jwt.JWT
	// Firebase        *firebases.FirebaseModel
	// Tokbox          *tokbox.Tokbox
	// GCS             googleCloud.IGCS
}

type Collection struct {
	Users user.UsersInterface
	// Otp                otp.OtpInterface
	// MAppKey            mAppKey.MAppKeyInterface
	// MConsent           consent.MConsentInterface
	// ProductCategories  ProductCategories.ProCInterface
	// Drugstores         drugstores.DrugstoresInterface
	// DsProductMasters   dsProductMasters.DsProductMastersinterface
	// DrugstoresFavorite drugstoresFavorites.DrugstoresFavoriteInterface
	// DrugCategories     drugCategories.DcCInterface
	// Pharmacist         pharmacist.Interface
	// Chat               chats.Interface
	// ChatMessage        messages.Interface
	// AppUpdate          appupdates.Interface
	// ProductMasters     productMasters.ProductMastersInterface
	// ProductFavorite    productFavorite.ProductFavoriteInterface
	// PharmcareOrder     pharmcareOrder.PhramcareInterface
	// MasterAddress      masterAddress.MasterAddressInterface
	// Opd                opd.OpdInterface
	// UserAddress        userAddress.UserAddressInterface
	// Cart               cart.CartInterface
	// FMembers           fMembers.FMembersInterface
	// Settings           settings.SettingsInterface
}

func GetConnect(ctx context.Context, k string) *Connection {
	log.Println("======================================================================================")
	log.Println("GetConnect:",ctx.Value(k).(Connection))
	log.Println("======================================================================================")
	data, err := base64.StdEncoding.DecodeString(k)
	if err != nil {
		log.Fatal("error:", err)
	}
	log.Println("data:", string(data))
	log.Println("ctx:", ctx.Value(k))
	
	if v, ok := ctx.Value(k).(Connection); ok {
		return &v
	}

	log.Panicln("Service can't create Connection or ctx not match")
	return nil
}

func GetCollection(ctx context.Context, k string) *Collection {
	log.Println("======================================================================================")
	log.Println("GetCollection:",ctx.Value(k).(Collection))
	log.Println("======================================================================================")
	data, err := base64.StdEncoding.DecodeString(k)
	if err != nil {
		log.Fatal("error:", err)
	}
	log.Println("data:", string(data))
	log.Println("ctx:", ctx.Value(k))

	if v, ok := ctx.Value(k).(Collection); ok {
		return &v
	}
	log.Panicln("Seveice can't create Collection or ctx not math")
	return nil
}
