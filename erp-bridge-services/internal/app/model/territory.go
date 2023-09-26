package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type Territory struct {
	ID             int64     `orm:"column(id)" json:"id"`
	Code           string    `orm:"column(code)" json:"code"`
	Description    string    `orm:"column(description)" json:"description"`
	RegionID       int64     `orm:"column(region_id)" json:"region_id"`
	SalespersonID  int64     `orm:"column(salesperson_id)" json:"salesperson_id"`
	CustomerTypeID int64     `orm:"column(customer_type_id)" json:"customer_type_id"`
	SubDistrictID  int64     `orm:"column(sub_district_id)" json:"sub_district_id"`
	CreatedAt      time.Time `orm:"column(created_at)" json:"created_at"`
	UpdatedAt      time.Time `orm:"column(updated_at)" json:"updated_at"`
}

func init() {
	orm.RegisterModel(new(Territory))
}

func (m *Territory) TableName() string {
	return "territory"
}
