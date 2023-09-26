package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type Permission struct {
	ID         int64         `orm:"column(id)" json:"id"`
	ParentID   int64         `orm:"column(parent_id)" json:"-"`
	Name       string        `orm:"column(name)" json:"name"`
	Value      string        `orm:"column(value)" json:"value"`
	CreatedAt  time.Time     `orm:"column(created_at)" json:"created_at"`
	UpdatedAt  time.Time     `orm:"column(updated_at)" json:"updated_at"`
	Child      []*Permission `orm:"-" json:"child,omitempty"`
	GrandChild []*Permission `orm:"-" json:"grand_child,omitempty"`
	Status     int8          `orm:"column(status)" json:"status"`
}

func init() {
	orm.RegisterModel(new(Permission))
}

func (m *Permission) TableName() string {
	return "permission"
}
