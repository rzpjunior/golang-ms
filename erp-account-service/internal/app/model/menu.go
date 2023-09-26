package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type Menu struct {
	ID           int64       `orm:"column(id)" json:"id"`
	ParentID     int64       `orm:"column(parent_id)" json:"parent_id"`
	Title        string      `orm:"column(title)" json:"title"`
	Url          string      `orm:"column(url)" json:"url"`
	Icon         string      `orm:"column(icon)" json:"icon"`
	PermissionID int64       `orm:"column(permission_id)" json:"permission_id"`
	Order        int         `orm:"column(order)" json:"order"`
	CreatedAt    time.Time   `orm:"column(created_at)" json:"created_at"`
	UpdatedAt    time.Time   `orm:"column(updated_at)" json:"updated_at"`
	Parent       *Menu       `orm:"-" json:"parent,omitempty"`
	Permission   *Permission `orm:"-" json:"permission,omitempty"`
	Child        []*Menu     `orm:"-" json:"child,omitempty"`
	Status       int8        `orm:"column(status)" json:"status"`
}

func init() {
	orm.RegisterModel(new(Menu))
}

func (m *Menu) TableName() string {
	return "menu"
}
