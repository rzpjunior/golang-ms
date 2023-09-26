package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type Region struct {
	ID          int64     `orm:"column(id)" json:"id"`
	Code        string    `orm:"column(code)" json:"code"`
	Description string    `orm:"column(description)" json:"description"`
	Status      int8      `orm:"column(status)" json:"status"`
	CreatedAt   time.Time `orm:"column(created_at)" json:"created_at"`
	UpdatedAt   time.Time `orm:"column(updated_at)" json:"updated_at"`
}

func init() {
	orm.RegisterModel(new(Region))
}

func (m *Region) TableName() string {
	return "region"
}
