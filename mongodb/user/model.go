package user

import (
	"github.com/khotchapan/KonLakRod-api/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Users struct {
	mongodb.Model  `bson:",inline"`
	FirstName      string `json:"firstName" bson:"first_name,omitempty"`
	LastName       string `json:"lastName" bson:"last_name,omitempty"`
	Image          string `json:"profileImage" bson:"user_image,omitempty"`
	Email          string `json:"email,omitempty" bson:"email,omitempty"`
	PhoneNumber    string `json:"phoneNumber,omitempty" bson:"phone_number,omitempty"`
	Birthday       string `json:"birthday" bson:"birthday,omitempty"`
	Username       string `json:"username" bson:"username,omitempty"`
	Password       string `json:"-" bson:"password,omitempty"`
	Activate       bool   `json:"activate" bson:"activate"`
	FacebookID     string `json:"-" bson:"facebook_id,omitempty"`
	FacebookActive bool   `json:"-" bson:"facebook_active"`
	GoogleID       string `json:"-" bson:"google_id,omitempty"`
	GoogleActive   bool   `json:"-" bson:"google_active,omitempty"`
	UserToken      string `json:"-" bson:"user_token,omitempty"`
	UserSex        string `json:"userSex" bson:"user_sex,omitempty"`
	//Address        []*primitive.ObjectID `json:"address" bson:"address,omitempty"`
	// AcceptConsent []*primitive.ObjectID `json:"acceptConsent" bson:"accept_consent,omitempty"`
	// HealthInfo    *HealthInfo            `json:"healthInfo" bson:"health_info,omitempty"`
}
type UsersResponse struct {
	// mongodb.Model  `bson:",inline"`
	// FirstName      string `json:"firstName" bson:"first_name,omitempty"`
	// LastName       string `json:"lastName" bson:"last_name,omitempty"`
	// Image          string `json:"profileImage" bson:"user_image,omitempty"`
	// Email          string `json:"email,omitempty" bson:"email,omitempty"`
	// PhoneNumber    string `json:"phoneNumber,omitempty" bson:"phone_number,omitempty"`
	// Birthday       string `json:"birthday" bson:"birthday,omitempty"`
	// Username       string `json:"username" bson:"username,omitempty"`
	// Password       string `json:"-" bson:"password,omitempty"`
	// Activate       bool   `json:"activate" bson:"activate"`
	// FacebookID     string `json:"-" bson:"facebook_id,omitempty"`
	// FacebookActive bool   `json:"-" bson:"facebook_active"`
	// GoogleID       string `json:"-" bson:"google_id,omitempty"`
	// GoogleActive   bool   `json:"-" bson:"google_active,omitempty"`
	// UserToken      string `json:"-" bson:"user_token,omitempty"`
	// UserSex        string `json:"userSex" bson:"user_sex,omitempty"`
	Id       primitive.ObjectID `json:"id,omitempty"`
    Name     string             `json:"name,omitempty" validate:"required"`
    Location string             `json:"location,omitempty" validate:"required"`
    Title    string             `json:"title,omitempty" validate:"required"`
}
type GetAllUsersForm struct {
	mongodb.PageQuery
	Name *string `query:"name"`
}