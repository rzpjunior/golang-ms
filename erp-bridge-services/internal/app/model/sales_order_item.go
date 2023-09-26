package model

import "git.edenfarm.id/edenlabs/edenlabs/orm"

type SalesOrderItem struct {
	ID            int64   `orm:"column(id)" json:"id"`
	SalesOrderId  int64   `orm:"column(sales_order_id)" json:"sales_order_id"`
	ItemId        int64   `orm:"column(item_id)" json:"item_id"`
	OrderQty      float64 `orm:"column(order_qty)" json:"order_qty"`
	DefaultPrice  float64 `orm:"column(default_price)" json:"default_price"`
	UnitPrice     float64 `orm:"column(unit_price)" json:"unit_price"`
	TaxableItem   int32   `orm:"column(taxable_item)" json:"taxable_item"`
	TaxPercentage float64 `orm:"column(tax_percentage)" json:"tax_percentage"`
	ShadowPrice   float64 `orm:"column(shadow_price)" json:"shadow_price"`
	Subtotal      float64 `orm:"column(subtotal)" json:"subtotal"`
	Weight        float64 `orm:"column(weight)" json:"weight"`
	Note          string  `orm:"column(note)" json:"note"`
}

func init() {
	orm.RegisterModel(new(SalesOrderItem))
}

func (m *SalesOrderItem) TableName() string {
	return "sales_order_item"
}
