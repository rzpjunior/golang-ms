package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type Uom struct {
	ID             int64     `orm:"column(id)" json:"id"`
	Code           string    `orm:"column(code)" json:"code"`
	Description    string    `orm:"column(description)" json:"description"`
	Status         int8      `orm:"column(status)" json:"status"`
	DecimalEnabled int8      `orm:"column(decimal_enabled)" json:"decimal_enabled"`
	CreatedAt      time.Time `orm:"column(created_at)" json:"created_at"`
	UpdatedAt      time.Time `orm:"column(updated_at)" json:"updated_at"`
}

func init() {
	orm.RegisterModel(new(Uom))
}

func (m *Uom) TableName() string {
	return "uom"
}
