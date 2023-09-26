package model

import (
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(DashboardOverview))
}

// DashboardTotalCharge model for Total Charge in Dashboard.
type DashboardOverview struct {
	ID               int64     `orm:"column(id);auto" json:"-,omitempty"`
	TotalTransaction float64   `orm:"column(total_transaction);null" json:"total_transaction"`
	SumTotalCharge 	 float64   `orm:"column(total_charge);null" json:"total_charge"`
	TotalTonnage 	 float64   `orm:"column(total_weight);null" json:"total_weight"`
	TopRevenue       []*toprevenue `orm:"-" json:"top_revenue"`
}

type toprevenue struct {
	ID    int64    `orm:"column(id);auto" json:"id"`
	Name  string   `orm:"column(name);null" json:"name"`
	Total float64  `orm:"column(total);null" json:"total"`
}