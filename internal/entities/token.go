package entities

import (
	"github.com/khotchapan/KonLakRod-api/internal/core/mongodb"
	"github.com/khotchapan/KonLakRod-api/internal/core/mongodb/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tokens struct {
	mongodb.Model `bson:",inline"`
	Token         string      `json:"token" bson:"token,omitempty"`
	RefreshToken  string      `json:"refresh_token" bson:"refresh_token,omitempty"`
	DeviceToken   string      `json:"device_token,omitempty" bson:"device_token,omitempty"`
	UserRefId        *primitive.ObjectID       `json:"user_ref_id" bson:"user_ref_id,omitempty"`
	User          *user.Users `json:"user" bson:"user,omitempty"`
}

type TokenResponse struct {
	AccessToken  *string `json:"access_token"`
	RefreshToken *string `json:"refresh_token"`
}
