package model

import (
	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

func init() {
	orm.RegisterModel(new(PaymentMethod))
}

// PaymentMethod : struct to hold payment term model data for database
type PaymentMethod struct {
	ID          int64  `orm:"column(id);auto" json:"-"`
	Code        string `orm:"column(code)" json:"code"`
	Name        string `orm:"column(name)" json:"name"`
	Note        string `orm:"column(note)" json:"note"`
	Status      int8   `orm:"column(status)" json:"status"`
	Publish     int8   `orm:"column(publish)" json:"publish"`
	Maintenance int8   `orm:"column(maintenance)" json:"maintenance"`
}

// TableName : set table name used by model
func (PaymentMethod) TableName() string {
	return "payment_method"
}
