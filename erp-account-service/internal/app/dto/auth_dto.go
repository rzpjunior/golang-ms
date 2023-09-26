package dto

type LoginRequest struct {
	Email    string `json:"email" valid:"required"`
	Password string `json:"password" valid:"required"`
	Timezone string
}

type LoginResponse struct {
	Token string `json:"token,omitempty"`
}

type UserPasswordRequest struct {
	Email    string `json:"email" valid:"required"`
	Password string `json:"password" valid:"required"`
	Timezone string
}

type UserPasswordResponse struct {
	UserResponse
}
