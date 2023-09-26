package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type RolePermission struct {
	ID           int64     `orm:"column(id)" json:"id"`
	RoleID       int64     `orm:"column(role_id)" json:"role_id,omitempty"`
	PermissionID int64     `orm:"column(permission_id)" json:"permission_id,omitempty"`
	CreatedAt    time.Time `orm:"column(created_at)" json:"created_at"`
	UpdatedAt    time.Time `orm:"column(updated_at)" json:"updated_at"`
	Status       int8      `orm:"column(status)" json:"status"`
}

func init() {
	orm.RegisterModel(new(RolePermission))
}

func (m *RolePermission) TableName() string {
	return "role_permission"
}
