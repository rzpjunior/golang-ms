package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type SalesPaymentTerm struct {
	ID          int64     `orm:"column(id)" json:"id"`
	Code        string    `orm:"column(code)" json:"code"`
	Description string    `orm:"column(description)" json:"description"`
	Status      int8      `orm:"column(status)" json:"status"`
	CreatedAt   time.Time `orm:"column(start_time)" json:"created_at"`
	UpdatedAt   time.Time `orm:"column(end_time)" json:"updated_at"`
}

func init() {
	orm.RegisterModel(new(SalesPaymentTerm))
}

func (m *SalesPaymentTerm) OrderType() string {
	return "term_payment_sls"
}
