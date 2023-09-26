package dto

import "time"

type LoginRequest struct {
	Email    string `json:"email" valid:"required"`
	Password string `json:"password" valid:"required"`
	FCM      string `json:"fcm_token"`
	Timezone string
}

type LoginResponse struct {
	Token string        `json:"token,omitempty"`
	User  *UserResponse `json:"user"`
}

type TokenClaim struct {
	UserId      int64
	Permissions []string
	ExpiresAt   time.Time
	Timezone    string
}
