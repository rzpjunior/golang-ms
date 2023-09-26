package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type PickingOrderAssign struct {
	ID                int64     `orm:"column(id)" json:"id"`
	PickingOrderId    int64     `orm:"column(picking_order_id)" json:"picking_order_id"`
	SiteID            string    `orm:"column(site_id)" json:"site_id"`
	SopNumber         string    `orm:"column(sop_number)" json:"sop_number"`
	DeliveryDate      time.Time `orm:"column(delivery_date)" json:"delivery_date"`
	WrtIdGP           string    `orm:"column(wrt_id_gp)" json:"wrt_id_gp"`
	CheckerIdGp       string    `orm:"column(checker_id_gp)" json:"checker_id"`
	CheckingStartTime time.Time `orm:"column(checking_start_time)" json:"checking_start_time"`
	CheckingEndTime   time.Time `orm:"column(checking_end_time)" json:"checking_end_time"`
	Status            int8      `orm:"column(status)" json:"status"`
}

func init() {
	orm.RegisterModel(new(PickingOrderAssign))
}

func (m *PickingOrderAssign) TableName() string {
	return "picking_order_assign"
}
