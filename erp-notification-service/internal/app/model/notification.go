package model

import "git.edenfarm.id/edenlabs/edenlabs/orm"

type Notification struct {
	ID      int64  `orm:"column(id)" json:"id"`
	Code    string `orm:"column(code)" json:"code"`
	Type    int8   `orm:"column(type)" json:"type"`
	Title   string `orm:"column(title)" json:"title"`
	Message string `orm:"column(message)" json:"message"`
	Status  int8   `orm:"column(status)" json:"status"`
}

func init() {
	orm.RegisterModel(new(Notification))
}

func (m *Notification) TableName() string {
	return "notification"
}
