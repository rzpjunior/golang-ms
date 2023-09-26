package model

import (
	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type RegionPolicy struct {
	ID                 int64  `orm:"column(id)" json:"-"`
	Region             string `orm:"column(region)" json:"region"`
	RegionID           int64  `orm:"column(region_id)" json:"region_id"`
	OrderTimeLimit     string `orm:"column(order_time_limit)" json:"order_time_limit"`
	MaxDayDeliveryDate int    `orm:"column(max_day_delivery_date)" json:"max_day_delivery_date"`
	WeeklyDayOff       int    `orm:"column(weekly_day_off)" json:"weekly_day_off"`
}

func init() {
	orm.RegisterModel(new(RegionPolicy))
}

func (m *RegionPolicy) TableName() string {
	return "region_policy"
}
