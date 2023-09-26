package dto

import "time"

type RoleResponse struct {
	ID            int64                 `json:"id,omitempty"`
	Code          string                `json:"code,omitempty"`
	Name          string                `json:"name,omitempty"`
	Division      *DivisionResponse     `json:"division,omitempty"`
	PermissionIDs []int64               `json:"permission_ids,omitempty"`
	Permissions   []*PermissionResponse `json:"permissions,omitempty"`
	CreatedAt     time.Time             `json:"created_at,omitempty"`
	UpdatedAt     time.Time             `json:"updated_at,omitempty"`
	Status        int8                  `json:"status,omitempty"`
	StatusConvert string                `json:"status_convert"`
	Note          string                `json:"note"`
}

type RoleRequestCreate struct {
	Name        string  `json:"name" valid:"required"`
	DivisionID  int64   `json:"division_id" valid:"required"`
	Permissions []int64 `json:"permissions,omitempty"`
	Note        string  `json:"note" valid:"lte:250"`
}

type RoleRequestUpdate struct {
	Name        string  `json:"name" valid:"required"`
	DivisionID  int64   `json:"division_id" valid:"required"`
	Permissions []int64 `json:"permissions,omitempty"`
	Note        string  `json:"note" valid:"lte:250"`
}
