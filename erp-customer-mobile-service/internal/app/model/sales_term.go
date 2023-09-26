package model

import "git.edenfarm.id/edenlabs/edenlabs/orm"

func init() {
	orm.RegisterModel(new(SalesTerm))
}

// SalesTerm : struct to hold payment term model data for database
type SalesTerm struct {
	ID          int64  `orm:"column(id);auto" json:"-"`
	Code        string `orm:"column(code);size(50);null" json:"code,omitempty"`
	Name        string `orm:"column(name);size(100);null" json:"name,omitempty"`
	DaysValue   int64  `orm:"column(days_value);null" json:"days_value,omitempty"`
	Note        string `orm:"column(note)" json:"note,omitempty"`
	Description string `orm:"column(description)" json:"description,omitempty"`
	Status      int8   `orm:"column(status);null" json:"status"`
}

// TableName : set table name used by model
func (SalesTerm) TableName() string {
	return "term_payment_sls"
}
