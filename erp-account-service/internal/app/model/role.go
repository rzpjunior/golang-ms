package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type Role struct {
	ID         int64     `orm:"column(id)" json:"id"`
	Code       string    `orm:"column(code)" json:"code"`
	Name       string    `orm:"column(name)" json:"name"`
	DivisionID int64     `orm:"column(division_id)" json:"division_id"`
	CreatedAt  time.Time `orm:"column(created_at)" json:"created_at"`
	UpdatedAt  time.Time `orm:"column(updated_at)" json:"updated_at"`
	Status     int8      `orm:"column(status)" json:"status"`
	Note       string    `orm:"column(note)" json:"note"`
}

func init() {
	orm.RegisterModel(new(Role))
}

func (m *Role) TableName() string {
	return "role"
}
