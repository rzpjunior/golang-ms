package model

import (
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(DashboardGraph))
}

// DashboardTotalCharge model for Total Charge in Dashboard.
type DashboardGraph struct {
	ID              int64    `orm:"column(id);auto" json:"-,omitempty"`
	Date  			string   `orm:"column(date);null" json:"date"`
	Day  			string   `orm:"column(day);null" json:"day"`
	TotalPrice		string   `orm:"column(total_price);null" json:"total_price"`
}
