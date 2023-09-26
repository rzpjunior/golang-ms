package model

import "git.edenfarm.id/edenlabs/edenlabs/orm"

type CourierVendor struct {
	ID     int64  `orm:"column(id)" json:"id"`
	Code   string `orm:"column(code)" json:"code"`
	Name   string `orm:"column(name)" json:"name"`
	Note   string `orm:"column(note)" json:"note"`
	Status int8   `orm:"column(status)" json:"status"`

	SiteID int64 `orm:"column(site_id)" json:"site_id"`
}

func init() {
	orm.RegisterModel(new(CourierVendor))
}

func (m *CourierVendor) TableName() string {
	return "courier_vendor"
}
