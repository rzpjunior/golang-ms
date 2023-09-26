package dto

type PurchaseOrderItemResponse struct {
	ID                 string                    `json:"id"`
	PurchaseOrderID    string                    `json:"purchase_order_id"`
	PurchasePlanItemID string                    `json:"purchase_plan_item_id"`
	PurchasePlanItem   *PurchasePlanItemResponse `json:"purchase_plan_item"`
	Item               *ItemResponse             `json:"item"`
	OrderQty           float64                   `json:"order_qty"`
	UnitPrice          float64                   `json:"unit_price"`
	TaxableItem        int32                     `json:"taxable_item"`
	IncludeTax         int32                     `json:"include_tax"`
	TaxPercentage      float64                   `json:"tax_percentage"`
	TaxAmount          float64                   `json:"tax_amount"`
	UnitPriceTax       float64                   `json:"unit_price_tax"`
	Subtotal           float64                   `json:"subtotal"`
	Weight             float64                   `json:"weight"`
	Note               string                    `json:"note"`
	PurchaseQty        float64                   `json:"purchase_qty"`
}

type PurchaseOrderItemRequestCreate struct {
	ItemID             string  `json:"item_id" valid:"required"`
	PurchasePlanItemID string  `json:"purchase_plan_item_id"`
	OrderQty           float64 `json:"qty" valid:"required"`
	UnitPrice          float64 `json:"unit_price" valid:"required"`
	Note               string  `json:"note" valid:"lte:500"`
	PurchaseQty        float64 `json:"purchase_qty"`
	IncludeTax         int8    `json:"include_tax"`
	TaxableItem        int8    `json:"taxable_item"`
	TaxPercentage      float64 `json:"tax_percentage"`
}

type PurchaseOrderItem struct {
	ID                int64   `orm:"column(id);auto" json:"-"`
	OrderQty          float64 `orm:"column(order_qty)" json:"order_qty"`
	UnitPrice         float64 `orm:"column(unit_price)" json:"unit_price"`
	Subtotal          float64 `orm:"column(subtotal)" json:"subtotal"`
	Weight            float64 `orm:"column(weight)" json:"weight"`
	Note              string  `orm:"column(note)" json:"note"`
	MarketPurchaseStr string  `orm:"column(market_purchase)" json:"market_purchase_str"`
	PurchaseQty       float64 `orm:"column(purchase_qty)" json:"purchase_qty"`
	TaxableItem       int8    `orm:"column(taxable_item)" json:"taxable_item"`
	IncludeTax        int8    `orm:"column(include_tax)" json:"include_tax"`
	TaxPercentage     float64 `orm:"column(tax_percentage)" json:"tax_percentage"`
	TaxAmount         float64 `orm:"column(tax_amount)" json:"tax_amount"`
	UnitPriceTax      float64 `orm:"column(unit_price_tax)" json:"unit_price_tax"`
	PurchaseOrder     string  `orm:"column(purchase_order_id);null;rel(fk)" json:"purchase_order"`
	Product           string  `orm:"column(product_id);null;rel(fk)" json:"product"`
}
