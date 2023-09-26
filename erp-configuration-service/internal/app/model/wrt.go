package model

import "git.edenfarm.id/edenlabs/edenlabs/orm"

type Wrt struct {
	ID    int64  `orm:"column(id)" json:"id"`
	WrtID string `orm:"column(wrt_id)" json:"wrt_id"`
	Type  int8   `orm:"column(type)" json:"type"`
	Note  string `orm:"column(note)" json:"note"`
}

func init() {
	orm.RegisterModel(new(Wrt))
}

func (m *Wrt) TableName() string {
	return "wrt"
}
