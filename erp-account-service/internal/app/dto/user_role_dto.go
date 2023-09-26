package dto

import "time"

type UserRoleResponse struct {
	ID        int64     `json:"-"`
	UserID    int64     `json:"user_id,omitempty"`
	RoleID    int64     `json:"role_id,omitempty"`
	MainRole  int8      `json:"main_role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Status    int8      `json:"status"`
}

type UserRoleByUserIdRequest struct {
	ID int64 `json:"id"`
}

type UserRoleByUserIdResponse struct {
	Roles []*UserRoleResponse `json:"roles"`
}
