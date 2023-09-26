package model

import "git.edenfarm.id/edenlabs/edenlabs/orm"

func init() {
	orm.RegisterModel(new(InvoiceTerm))
}

// Invoice Term : struct to hold payment term model data for database
type InvoiceTerm struct {
	ID     int64  `orm:"column(id);auto" json:"-"`
	Code   string `orm:"column(code);size(50);null" json:"code"`
	Name   string `orm:"column(name);size(100);null" json:"name"`
	NameID string `orm:"column(name_id);size(100);null" json:"name_id"`
	Note   string `orm:"column(note)" json:"note"`
	Status int8   `orm:"column(status);null" json:"status"`
}

// TableName : set table name used by model
func (InvoiceTerm) TableName() string {
	return "term_invoice_sls"
}
