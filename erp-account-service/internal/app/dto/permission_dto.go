package dto

import "time"

type PermissionResponse struct {
	ID            int64                 `json:"id,omitempty"`
	ParentID      int64                 `json:"parent_id,omitempty"`
	Name          string                `json:"name,omitempty"`
	Value         string                `json:"value,omitempty"`
	CreatedAt     time.Time             `json:"created_at,omitempty"`
	UpdatedAt     time.Time             `json:"updated_at,omitempty"`
	Child         []*PermissionResponse `json:"child,omitempty"`
	GrandChild    []*PermissionResponse `json:"grand_child,omitempty"`
	Status        int8                  `json:"status,omitempty"`
	StatusConvert string                `json:"status_convert"`
}
