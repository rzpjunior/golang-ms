package model

import "git.edenfarm.id/edenlabs/edenlabs/orm"

type PurchasePlanItem struct {
	ID              int64   `orm:"column(id)" json:"id"`
	PurchasePlanID  int64   `orm:"column(purchase_plan_id)" json:"purchase_plan_id"`
	ItemID          int64   `orm:"column(item_id)" json:"item_id"`
	PurchasePlanQty float64 `orm:"column(purchase_plan_qty)" json:"purchase_plan_qty"`
	PurchaseQty     float64 `orm:"column(purchase_qty)" json:"purchase_qty"`
	UnitPrice       float64 `orm:"column(unit_price)" json:"unit_price"`
	Subtotal        float64 `orm:"column(subtotal)" json:"subtotal"`
	Weight          float64 `orm:"column(weight)" json:"weight"`
}

func init() {
	orm.RegisterModel(new(PurchasePlanItem))
}

func (m *PurchasePlanItem) TableName() string {
	return "purchase_plan_item"
}
