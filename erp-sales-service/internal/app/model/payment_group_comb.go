package model

import (
	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

func init() {
	orm.RegisterModel(new(PaymentGroupComb))
}

// PaymentGroupComb : struct to hold payment term model data for database
type PaymentGroupComb struct {
	ID              int64  `orm:"column(id);auto" json:"-"`
	PaymentGroupSls string `orm:"column(payment_group_sls);" json:"payment_group_sls"`
	TermPaymentSls  string `orm:"column(term_payment_sls);" json:"term_payment_sls"`
}

// TableName : set table name used by model
func (PaymentGroupComb) TableName() string {
	return "payment_group_comb"
}
