package entities

import "github.com/khotchapan/KonLakRod-api/internal/core/mongodb/user"

// Token token
type Token struct {
	Model
	Token        string      `json:"token"`
	RefreshToken string      `json:"refresh_token"`
	DeviceToken  string      `json:"device_token,omitempty"`
	UserID       string      `json:"user_id"`
	User         *user.Users `json:"user"`
}
