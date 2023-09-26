package model

import (
	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type Koli struct {
	Id     int64  `orm:"column(id)" json:"id"`
	Code   string `orm:"column(code)" json:"code"`
	Value  string `orm:"column(value)" json:"value"`
	Name   string `orm:"column(name)" json:"name"`
	Note   string `orm:"column(note)" json:"note"`
	Status int8   `orm:"column(status)" json:"status"`
}

func init() {
	orm.RegisterModel(new(Koli))
}

func (m *Koli) TableName() string {
	return "koli"
}
