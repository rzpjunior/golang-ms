package model

import "git.edenfarm.id/edenlabs/edenlabs/orm"

type PurchaseOrderItem struct {
	ID                 int64   `orm:"column(id)" json:"id"`
	PurchaseOrderID    int64   `orm:"column(purchase_order_id)" json:"purchase_order_id"`
	PurchasePlanItemID int64   `orm:"column(purchase_plan_item_id)" json:"purchase_plan_item_id"`
	ItemID             int64   `orm:"column(item_id)" json:"item_id"`
	OrderQty           float64 `orm:"column(order_qty)" json:"order_qty"`
	UnitPrice          float64 `orm:"column(unit_price)" json:"unit_price"`
	TaxableItem        int32   `orm:"column(taxable_item)" json:"taxable_item"`
	IncludeTax         int32   `orm:"column(include_tax)" json:"include_tax"`
	TaxPercentage      float64 `orm:"column(tax_percentage)" json:"tax_percentage"`
	TaxAmount          float64 `orm:"column(tax_amount)" json:"tax_amount"`
	UnitPriceTax       float64 `orm:"column(unit_price_tax)" json:"unit_price_tax"`
	Subtotal           float64 `orm:"column(subtotal)" json:"subtotal"`
	Weight             float64 `orm:"column(weight)" json:"weight"`
	Note               string  `orm:"column(note)" json:"note"`
	PurchaseQty        float64 `orm:"column(purchase_qty)" json:"purchase_qty"`
}

func init() {
	orm.RegisterModel(new(PurchaseOrderItem))
}

func (m *PurchaseOrderItem) TableName() string {
	return "purchase_order_item"
}
