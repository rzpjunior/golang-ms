package model

import "git.edenfarm.id/edenlabs/edenlabs/orm"

func init() {
	orm.RegisterModel(new(PriceSet))
}

// PriceSet : struct to hold price set model data for database
type PriceSet struct {
	ID     int64  `orm:"column(id);auto" json:"-"`
	Code   string `orm:"column(code)" json:"code"`
	Name   string `orm:"column(name)" json:"name"`
	Note   string `orm:"column(note)" json:"note"`
	Status int8   `orm:"column(status)" json:"status"`
}
