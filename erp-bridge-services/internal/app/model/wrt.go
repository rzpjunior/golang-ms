package model

import "git.edenfarm.id/edenlabs/edenlabs/orm"

type Wrt struct {
	ID        int64  `orm:"column(id)" json:"id"`
	RegionID  int64  `orm:"column(region_id)" json:"region_id"`
	Code      string `orm:"column(code)" json:"code"`
	StartTime string `orm:"column(start_time)" json:"start_time"`
	EndTime   string `orm:"column(end_time)" json:"end_time"`
}

func init() {
	orm.RegisterModel(new(Wrt))
}

func (m *Wrt) TableName() string {
	return "wrt"
}
