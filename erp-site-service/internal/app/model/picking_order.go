package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type PickingOrder struct {
	Id        int64     `orm:"column(id)" json:"id"`
	DocNumber string    `orm:"column(doc_number)" json:"doc_number"`
	PickerId  string    `orm:"column(picker_id)" json:"picker_id"`
	StartTime time.Time `orm:"column(start_time)" json:"start_time"`
	EndTime   time.Time `orm:"column(end_time)" json:"end_time"`
	Status    int8      `orm:"column(status)" json:"status"`
}

func init() {
	orm.RegisterModel(new(PickingOrder))
}

func (m *PickingOrder) TableName() string {
	return "picking_order"
}
