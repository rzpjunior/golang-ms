package model

import (
	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type CustomerPointSummary struct {
	ID            int64   `orm:"column(id);auto" json:"id"`
	CustomerID    int64   `orm:"column(customer_id)" json:"customer_id"`
	EarnedPoint   float64 `orm:"column(earned_point);null" json:"earned_point"`
	RedeemedPoint float64 `orm:"column(redeemed_point);null" json:"redeemed_point"`
	SummaryDate   string  `orm:"column(summary_date);null" json:"summary_date"`
}

func init() {
	orm.RegisterModel(new(CustomerPointSummary))
}

func (m *CustomerPointSummary) TableName() string {
	return "customer_point_summary"
}
