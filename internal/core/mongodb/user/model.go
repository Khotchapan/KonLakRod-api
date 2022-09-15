package user

import (
	"github.com/khotchapan/KonLakRod-api/internal/core/mongodb"
)

type Users struct {
	mongodb.Model  `bson:",inline"`
	FirstName      string   `json:"firstName" bson:"first_name,omitempty"`
	LastName       string   `json:"lastName" bson:"last_name,omitempty"`
	Image          string   `json:"image" bson:"image,omitempty"`
	Email          string   `json:"email,omitempty" bson:"email,omitempty"`
	PhoneNumber    string   `json:"phoneNumber,omitempty" bson:"phone_number,omitempty"`
	Birthday       string   `json:"birthday" bson:"birthday,omitempty"`
	Username       string   `json:"username" bson:"username,omitempty"`
	PasswordHash   string   `json:"password_hash" bson:"password_hash,omitempty"`
	Roles          []string `json:"roles,omitempty" bson:"roles,omitempty"`
	Activate       bool     `json:"activate" bson:"activate,omitempty"`
	FacebookID     string   `json:"-" bson:"facebook_id,omitempty"`
	FacebookActive bool     `json:"-" bson:"facebook_active,omitempty"`
	GoogleID       string   `json:"-" bson:"google_id,omitempty"`
	GoogleActive   bool     `json:"-" bson:"google_active,omitempty"`
	UserToken      string   `json:"-" bson:"user_token,omitempty"`
	UserSex        string   `json:"userSex" bson:"user_sex,omitempty"`
	//Address        []*primitive.ObjectID `json:"address" bson:"address,omitempty"`
	// AcceptConsent []*primitive.ObjectID `json:"acceptConsent" bson:"accept_consent,omitempty"`
	// HealthInfo    *HealthInfo            `json:"healthInfo" bson:"health_info,omitempty"`
}
type UsersResponse struct {
	mongodb.Model  `bson:",inline"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Image          string `json:"profile_image"`
	Email          string `json:"email,omitempty"`
	PhoneNumber    string `json:"phone_number,omitempty"`
	Birthday       string `json:"birthday"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	Activate       bool   `json:"activate"`
	FacebookID     string `json:"facebook_id"`
	FacebookActive bool   `json:"facebook_active"`
	GoogleID       string `json:"google_id"`
	GoogleActive   bool   `json:"google_active"`
	UserToken      string `json:"user_token"`
	UserSex        string `json:"user_sex"`
}
type GetAllUsersForm struct {
	mongodb.PageQuery
	Name *string `query:"name"`
}
