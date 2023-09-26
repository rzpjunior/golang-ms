package model

import (
	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type PickingOrderItem struct {
	Id                   int64   `orm:"column(id)" json:"id"`
	IdGp                 int64   `orm:"column(picking_order_item_id_gp)" json:"-"`
	PickingOrderAssignId int64   `orm:"column(picking_order_assign_id)" json:"picking_order_assign_id"`
	ItemNumber           string  `orm:"column(item_number)" json:"item_number"`
	OrderQuantity        float64 `orm:"column(order_qty)" json:"order_quantity"`
	PickQuantity         float64 `orm:"column(pick_qty)" json:"pick_quantity"`
	CheckQuantity        float64 `orm:"column(check_qty)" json:"check_quantity"`
	ExcessQuantity       float64 `orm:"column(excess_qty)" json:"excess_quantity"`
	UnfulfillNote        string  `orm:"column(unfulfill_note)" json:"unfulfill_note"`
	Status               int8    `orm:"column(status)" json:"status"`
}

func init() {
	orm.RegisterModel(new(PickingOrderItem))
}

func (m *PickingOrderItem) TableName() string {
	return "picking_order_item"
}
