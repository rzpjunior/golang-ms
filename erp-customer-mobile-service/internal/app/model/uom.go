package model

import "git.edenfarm.id/edenlabs/edenlabs/orm"

func init() {
	orm.RegisterModel(new(Uom))
}

// Uom : struct to hold uom model data for database
type Uom struct {
	ID             int64  `orm:"column(id);auto" json:"-"`
	Code           string `orm:"column(code)" json:"code"`
	Name           string `orm:"column(name)" json:"name"`
	DecimalEnabled int8   `orm:"column(decimal_enabled)" json:"decimal_enabled"`
	Note           string `orm:"column(note)" json:"note"`
	Status         int8   `orm:"column(status)" json:"status"`
}
