package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type UserRole struct {
	ID        int64     `orm:"column(id)" json:"id"`
	UserID    int64     `orm:"column(user_id)" json:"user_id,omitempty"`
	RoleID    int64     `orm:"column(role_id)" json:"role_id,omitempty"`
	MainRole  int8      `orm:"column(main_role)" json:"main_role"`
	CreatedAt time.Time `orm:"column(created_at)" json:"created_at"`
	UpdatedAt time.Time `orm:"column(updated_at)" json:"updated_at"`
	Status    int8      `orm:"column(status)" json:"status"`
}

func init() {
	orm.RegisterModel(new(UserRole))
}

func (m *UserRole) TableName() string {
	return "user_role"
}
